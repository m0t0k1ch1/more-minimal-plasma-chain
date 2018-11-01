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
		cmdTx,
	}

	if err := app.Run(os.Args); err != nil {
		printlnJSON(map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}
}
