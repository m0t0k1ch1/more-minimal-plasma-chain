package main

import "github.com/urfave/cli"

var cmdTxOut = cli.Command{
	Name:  "txout",
	Usage: "commands for txout",
	Subcommands: []cli.Command{
		cmdTxOutExit,
		cmdTxOutRestore,
	},
}
