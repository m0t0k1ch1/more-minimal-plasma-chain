package app

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

const (
	DefaultMempoolSize = 1000
)

type HandlerFunc func(*Context) error

type ChildChain struct {
	e          *echo.Echo
	config     *Config
	operator   *types.Account
	blockchain *core.Blockchain
}

func NewChildChain(conf *Config) (*ChildChain, error) {
	privKey, err := crypto.HexToECDSA(conf.Operator.PrivateKey)
	if err != nil {
		return nil, err
	}

	cc := &ChildChain{
		e:          echo.New(),
		config:     conf,
		operator:   types.NewAccount(privKey),
		blockchain: core.NewBlockchain(),
	}

	cc.e.Use(middleware.Logger())
	cc.e.Use(middleware.Recover())
	cc.e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&Context{c})
		}
	})
	cc.e.HTTPErrorHandler = cc.httpErrorHandler

	cc.GET("/ping", cc.PingHandler)
	cc.POST("/blocks", cc.PostBlockHandler)
	cc.GET("/blocks/:blkNum", cc.GetBlockHandler)
	cc.GET("/blocks/:blkNum/txes/:txIndex", cc.GetBlockTxHandler)
	cc.GET("/blocks/:blkNum/txes/:txIndex/proof", cc.GetBlockTxProofHandler)
	cc.PUT("/blocks/:blkNum/txes/:txIndex/inputs/:iIndex", cc.PutBlockTxInHandler)
	cc.POST("/txes", cc.PostTxHandler)

	return cc, nil
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
	cc.e.Add(method, path, func(c echo.Context) error {
		return h(NewContext(c))
	})
}

func (cc *ChildChain) Logger() echo.Logger {
	return cc.e.Logger
}

func (cc *ChildChain) Start() error {
	return cc.e.Start(fmt.Sprintf(":%d", cc.config.Port))
}

func (cc *ChildChain) httpErrorHandler(err error, c echo.Context) {
	cc.e.Logger.Error(err)

	code := http.StatusInternalServerError
	msg := http.StatusText(code)

	if httpErr, ok := err.(*echo.HTTPError); ok {
		code = httpErr.Code
		msg = fmt.Sprintf("%v", httpErr.Message)
	}

	appErr := NewError(code, msg)

	if err := c.JSON(appErr.Code, NewErrorResponse(appErr)); err != nil {
		cc.e.Logger.Error(err)
	}
}
