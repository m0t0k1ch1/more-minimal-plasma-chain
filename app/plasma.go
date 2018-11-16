package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type HandlerFunc func(*Context) error

type Plasma struct {
	config            Config
	server            *echo.Echo
	db                *DB
	operator          *types.Account
	rootChain         *core.RootChain
	childChain        *core.ChildChain
	heartbeater       *Heartbeater
	heartbeatInterval time.Duration
}

func NewPlasma(conf Config) (*Plasma, error) {
	p := &Plasma{
		config: conf,
	}

	p.initServer()
	if err := p.initDB(); err != nil {
		return nil, err
	}
	if err := p.initRootChain(); err != nil {
		return nil, err
	}
	if err := p.initOperator(); err != nil {
		return nil, err
	}
	p.initChildChain()

	if conf.Heartbeat.IsEnabled {
		if err := p.initHeartbeater(); err != nil {
			return nil, err
		}
		if err := p.initHeartbeatInterval(); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *Plasma) initServer() {
	p.server = echo.New()
	p.server.Use(middleware.Logger())
	p.server.Use(middleware.Recover())
	p.server.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&Context{c})
		}
	})
	p.server.HTTPErrorHandler = p.httpErrorHandler
	p.server.Logger.SetLevel(log.INFO)
	p.initRoutes()
}

func (p *Plasma) initDB() error {
	db, err := NewDB(p.config.DB)
	if err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p *Plasma) initRoutes() {
	p.GET("/ping", p.PingHandler)
	p.GET("/addresses/:address/utxos", p.GetAddressUTXOsHandler)
	p.POST("/blocks", p.PostBlockHandler)
	p.GET("/blocks/:blkNum", p.GetBlockHandler)
	p.POST("/txes", p.PostTxHandler)
	p.GET("/txes/:txPos", p.GetTxHandler)
	p.GET("/txes/:txPos/proof", p.GetTxProofHandler)
	p.PUT("/txins/:txInPos", p.PutTxInHandler)
	p.POST("/deposits", p.PostDepositHandler)
}

func (p *Plasma) initRootChain() error {
	rc, err := core.NewRootChain(p.config.RootChain)
	if err != nil {
		return err
	}
	p.rootChain = rc
	return nil
}

func (p *Plasma) initOperator() error {
	privKey, err := p.config.Operator.PrivateKey()
	if err != nil {
		return err
	}
	p.operator = types.NewAccount(privKey)
	return nil
}

func (p *Plasma) initChildChain() error {
	return p.db.Update(func(txn *badger.Txn) error {
		cc, err := core.NewChildChain(txn)
		if err != nil {
			return err
		}
		p.childChain = cc
		return nil
	})
}

func (p *Plasma) initHeartbeater() error {
	heartbeater, err := NewHeartbeater(p.rootChain.Ping)
	if err != nil {
		return err
	}
	p.heartbeater = heartbeater
	return nil
}

func (p *Plasma) initHeartbeatInterval() error {
	interval, err := p.config.Heartbeat.Interval()
	if err != nil {
		return err
	}
	p.heartbeatInterval = interval
	return nil
}

func (p *Plasma) GET(path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	p.Add(http.MethodGet, path, h, m...)
}

func (p *Plasma) POST(path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	p.Add(http.MethodPost, path, h, m...)
}

func (p *Plasma) PUT(path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	p.Add(http.MethodPut, path, h, m...)
}

func (p *Plasma) Add(method, path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	p.server.Add(method, path, func(c echo.Context) error {
		return h(NewContext(c))
	})
}

func (p *Plasma) Logger() echo.Logger {
	return p.server.Logger
}

func (p *Plasma) Start() error {
	// watch DepositCreated events
	if err := p.watchDepositCreated(); err != nil {
		return err
	}

	// watch ExitStarted events
	if err := p.watchExitStarted(); err != nil {
		return err
	}

	if p.config.Heartbeat.IsEnabled {
		// keep WebSocket connection alive
		if err := p.heartbeat(); err != nil {
			return err
		}
	}

	// start HTTP server
	return p.server.Start(fmt.Sprintf(":%d", p.config.Port))
}

func (p *Plasma) Shutdown(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

func (p *Plasma) Finalize() {
	p.db.Close()

	if p.config.Heartbeat.IsEnabled {
		p.heartbeater.Stop()
	}
}

func (p *Plasma) watchDepositCreated() error {
	sink := make(chan *core.RootChainDepositCreated)
	if _, err := p.rootChain.WatchDepositCreated(context.Background(), sink); err != nil {
		return err
	}

	go func() {
		for log := range sink {
			if err := p.db.Update(func(txn *badger.Txn) error {
				newBlkNum, err := p.childChain.AddDepositBlock(txn, log.Owner, log.Amount.Uint64(), p.operator)
				if err != nil {
					p.Logger().Error(err)
				} else {
					p.Logger().Infof(
						"[DEPOSIT] owner: %s, amount: %d, blkNum: %d, txPos: %d",
						utils.AddressToHex(log.Owner),
						log.Amount,
						newBlkNum,
						types.NewTxPosition(newBlkNum, 0),
					)
				}
				return nil
			}); err != nil {
				p.Logger().Error(err)
			}
		}
	}()

	return nil
}

func (p *Plasma) watchExitStarted() error {
	sink := make(chan *core.RootChainExitStarted)
	if _, err := p.rootChain.WatchExitStarted(context.Background(), sink); err != nil {
		return err
	}

	go func() {
		for log := range sink {
			if err := p.db.Update(func(txn *badger.Txn) error {
				txOutPos := types.Position(log.UtxoPosition.Uint64())
				if err := p.childChain.ExitTxOut(txn, txOutPos); err != nil {
					p.Logger().Error(err)
				} else {
					p.Logger().Infof(
						"[EXIT] owner: %s, amount: %d, txOutPos: %d",
						utils.AddressToHex(log.Owner),
						log.Amount,
						txOutPos,
					)
				}
				return nil
			}); err != nil {
				p.Logger().Error(err)
			}
		}
	}()

	return nil
}

func (p *Plasma) heartbeat() error {
	go func() {
		for {
			ok, err := p.heartbeater.Beat()
			if err != nil {
				p.Logger().Error(err)
			}
			if !ok {
				return
			}

			time.Sleep(p.heartbeatInterval)
		}
	}()

	return nil
}

func (p *Plasma) httpErrorHandler(err error, c echo.Context) {
	p.Logger().Error(err)

	code := http.StatusInternalServerError
	msg := http.StatusText(code)

	if httpErr, ok := err.(*echo.HTTPError); ok {
		code = httpErr.Code
		msg = fmt.Sprintf("%v", httpErr.Message)
	}

	appErr := NewError(code, msg)

	if err := c.JSON(appErr.Code, NewErrorResponse(appErr)); err != nil {
		p.Logger().Error(err)
	}
}
