package main

import (
	"github.com/urfave/cli"
)

var cmdTx = cli.Command{
	Name:  "tx",
	Usage: "commands for tx",
	Subcommands: []cli.Command{
		cmdTxConfirm,
		cmdTxCreate,
		cmdTxGet,
		cmdTxPost,
		cmdTxProof,
		cmdTxSign,
	},
}
