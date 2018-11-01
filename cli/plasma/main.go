package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		CmdBlock,
		CmdTx,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
