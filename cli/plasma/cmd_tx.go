package main

import (
	"github.com/urfave/cli"
)

var CmdTx = cli.Command{
	Name:  "tx",
	Usage: "commands for tx",
	Subcommands: []cli.Command{
		CmdTxCreate,
		CmdTxGet,
		CmdTxPost,
		CmdTxSign,
	},
}
