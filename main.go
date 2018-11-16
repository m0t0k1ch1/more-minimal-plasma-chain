package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/app"
)

const (
	DefaultConfigPath = "config.json"
)

func loadConfig(path string) (app.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return app.Config{}, err
	}

	var conf app.Config
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		return app.Config{}, err
	}

	return conf, nil
}

func main() {
	var confPath = flag.String("conf", DefaultConfigPath, "path to your config.json")
	flag.Parse()

	conf, err := loadConfig(*confPath)
	if err != nil {
		panic(err)
	}

	p, err := app.NewPlasma(conf)
	if err != nil {
		panic(err)
	}

	done := make(chan struct{}, 0)
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGTERM)
		<-sigterm

		p.Finalize()
		if err := p.Shutdown(context.Background()); err != nil {
			p.Logger().Fatal(err)
		}
		close(done)
	}()
	if err := p.Start(); err != nil {
		p.Logger().Info(err)
	}
	<-done
}
