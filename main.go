package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	DefaultConfigPath = "config.json"
)

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func main() {
	var confPath = flag.String("conf", DefaultConfigPath, "path to your config.json")
	flag.Parse()

	conf, err := loadConfig(*confPath)
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ping", PingHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Port)))
}
