package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func PingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
