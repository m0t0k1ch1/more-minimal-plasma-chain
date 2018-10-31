package main

import (
	"log"
	"os"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/cli/plasma/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		cmd.CmdTx,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
