package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/models"
)

const (
	DefaultMempoolSize = 1000
)

type HandlerFunc func(*Context) error

type ChildChain struct {
	e          *echo.Echo
	config     *Config
	blockchain *models.Blockchain
	mempool    *models.Mempool
}

func NewChildChain(conf *Config) *ChildChain {
	cc := &ChildChain{
		e:          echo.New(),
		config:     conf,
		blockchain: models.NewBlockchain(),
		mempool:    models.NewMempool(DefaultMempoolSize),
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
	cc.POST("/depositBlocks", cc.PostDepositBlockHandler)

	return cc
}

func (cc *ChildChain) GET(path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	cc.Add(http.MethodGet, path, h, m...)
}

func (cc *ChildChain) POST(path string, h HandlerFunc, m ...echo.MiddlewareFunc) {
	cc.Add(http.MethodPost, path, h, m...)
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
