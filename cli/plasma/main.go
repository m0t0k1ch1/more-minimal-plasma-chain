package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		cmdBlock,
		cmdChain,
		cmdDeposit,
		cmdExit,
		cmdTx,
	}

	if err := app.Run(os.Args); err != nil {
		printlnJSON(map[string]string{
			"error": err.Error(),
		})
		os.Exit(1)
	}
}
