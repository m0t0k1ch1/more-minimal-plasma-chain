package main

import "github.com/urfave/cli"

var cmdTxIn = cli.Command{
	Name:  "txin",
	Usage: "commands for tx input",
	Subcommands: []cli.Command{
		cmdTxInConfirm,
	},
}
