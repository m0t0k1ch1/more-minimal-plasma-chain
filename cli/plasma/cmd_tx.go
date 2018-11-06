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
		cmdTxDecode,
		cmdTxGet,
		cmdTxIndex,
		cmdTxPost,
		cmdTxProof,
		cmdTxSign,
	},
}
