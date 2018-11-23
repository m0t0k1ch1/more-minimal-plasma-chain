package main

import (
	"encoding/json"
	"os"

	"github.com/urfave/cli"
)

var (
	conf Config
)

func loadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	return json.NewDecoder(file).Decode(&conf)
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		cmdAddress,
		cmdBlock,
		cmdDeploy,
		cmdDeposit,
		cmdExit,
		cmdTx,
		cmdTxIn,
		cmdTxOut,
	}
	app.Flags = []cli.Flag{
		confFlag,
	}
	app.Before = func(c *cli.Context) error {
		if err := loadConfig(getGlobalString(c, confFlag)); err != nil {
			exit(err)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		exit(err)
	}
}
