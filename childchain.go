package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type ChildChain struct {
	e      *echo.Echo
	config *Config
}

func NewChildChain(conf *Config) *ChildChain {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ping", PingHandler)

	return &ChildChain{
		e:      e,
		config: conf,
	}
}

func (cc *ChildChain) Logger() echo.Logger {
	return cc.e.Logger
}

func (cc *ChildChain) Start() error {
	return cc.e.Start(fmt.Sprintf(":%d", cc.config.Port))
}
