package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/models"
)

const (
	DefaultBlockNumber = 1
)

type HandlerFunc func(*Context) error

type ChildChain struct {
	e                  *echo.Echo
	config             *Config
	blocks             map[uint64]*models.Block
	currentBlockNumber uint64
}

func NewChildChain(conf *Config) *ChildChain {
	cc := &ChildChain{
		e:                  echo.New(),
		config:             conf,
		blocks:             map[uint64]*models.Block{},
		currentBlockNumber: DefaultBlockNumber,
	}

	cc.e.Use(middleware.Logger())
	cc.e.Use(middleware.Recover())
	cc.e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&Context{c})
		}
	})
	cc.GET("/ping", cc.PingHandler)
	cc.POST("/blocks", cc.PostBlockHandler)
	cc.GET("/blocks/:num", cc.GetBlockHandler)

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