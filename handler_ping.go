package main

import (
	"net/http"
)

func (cc *ChildChain) PingHandler(c *Context) error {
	return c.String(http.StatusOK, "pong")
}
