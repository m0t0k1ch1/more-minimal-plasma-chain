package main

import (
	"github.com/urfave/cli"
)

var cmdTx = cli.Command{
	Name:  "tx",
	Usage: "commands for tx",
	Subcommands: []cli.Command{
		cmdTxDecode,
		cmdTxGet,
		cmdTxPost,
		cmdTxProof,
	},
}
