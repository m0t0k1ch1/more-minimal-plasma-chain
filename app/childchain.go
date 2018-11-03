package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/contract"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type HandlerFunc func(*Context) error

type ChildChain struct {
	config     *Config
	server     *echo.Echo
	rootChain  *RootChain
	operator   *types.Account
	blockchain *core.Blockchain
}

func NewChildChain(conf *Config) (*ChildChain, error) {
	cc := &ChildChain{
		config: conf,
	}

	cc.initServer()
	if err := cc.initRootChain(); err != nil {
		return nil, err
	}
	if err := cc.initOperator(); err != nil {
		return nil, err
	}
	cc.initBlockchain()

	return cc, nil
}

func (cc *ChildChain) initServer() {
	cc.server = echo.New()
	cc.server.Use(middleware.Logger())
	cc.server.Use(middleware.Recover())
	cc.server.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&Context{c})
		}
	})
	cc.server.HTTPErrorHandler = cc.httpErrorHandler
	cc.server.Logger.SetLevel(log.INFO)
	cc.initRoutes()
}

func (cc *ChildChain) initRoutes() {
	cc.GET("/ping", cc.PingHandler)
	cc.GET("/chain/:blkNum", cc.GetChainHandler)
	cc.POST("/blocks", cc.PostBlockHandler)
	cc.GET("/blocks/:blkHash", cc.GetBlockHandler)
	cc.POST("/txes", cc.PostTxHandler)
	cc.GET("/txes/:txHash", cc.GetTxHandler)
	cc.GET("/txes/:txHash/proof", cc.GetTxProofHandler)
	cc.PUT("/txes/:txHash", cc.PutTxHandler)
}

func (cc *ChildChain) initRootChain() error {
	rc, err := NewRootChain(cc.config.RootChain)
	if err != nil {
		return err
	}
	cc.rootChain = rc
	return nil
}

func (cc *ChildChain) initOperator() error {
	privKey, err := crypto.HexToECDSA(cc.config.Operator.PrivateKey)
	if err != nil {
		return err
	}
	cc.operator = types.NewAccount(privKey)
	return nil
}

func (cc *ChildChain) initBlockchain() {
	cc.blockchain = core.NewBlockchain()
}

func (cc *ChildChain) GET(path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	cc.Add(http.MethodGet, path, h, m...)
}

func (cc *ChildChain) POST(path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	cc.Add(http.MethodPost, path, h, m...)
}

func (cc *ChildChain) PUT(path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	cc.Add(http.MethodPut, path, h, m...)
}

func (cc *ChildChain) Add(method, path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	cc.server.Add(method, path, func(c echo.Context) error {
		return h(NewContext(c))
	})
}

func (cc *ChildChain) Logger() echo.Logger {
	return cc.server.Logger
}

func (cc *ChildChain) Start() error {
	if err := cc.watchRootChain(); err != nil {
		return err
	}

	return cc.server.Start(fmt.Sprintf(":%d", cc.config.Port))
}

func (cc *ChildChain) watchRootChain() error {
	sink := make(chan *contract.RootChainDepositCreated)
	sub, err := cc.rootChain.WatchDepositCreated(context.Background(), sink)
	if err != nil {
		return err
	}

	go func() {
		defer sub.Unsubscribe()
		for log := range sink {
			blkHash, err := cc.blockchain.AddDepositBlock(log.Owner, log.Amount.Uint64(), cc.operator)
			if err != nil {
				cc.Logger().Error(err)
			} else {
				cc.Logger().Infof(
					"[DEPOSIT] blkhash: %s, owner: %s: amount: %d",
					utils.EncodeToHex(blkHash.Bytes()),
					utils.EncodeToHex(log.Owner.Bytes()),
					log.Amount.Uint64(),
				)
			}
		}
	}()

	return nil
}

func (cc *ChildChain) httpErrorHandler(err error, c echo.Context) {
	cc.Logger().Error(err)

	code := http.StatusInternalServerError
	msg := http.StatusText(code)

	if httpErr, ok := err.(*echo.HTTPError); ok {
		code = httpErr.Code
		msg = fmt.Sprintf("%v", httpErr.Message)
	}

	appErr := NewError(code, msg)

	if err := c.JSON(appErr.Code, NewErrorResponse(appErr)); err != nil {
		cc.Logger().Error(err)
	}
}
