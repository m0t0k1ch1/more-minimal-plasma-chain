package main

import "github.com/urfave/cli"

var cmdChain = cli.Command{
	Name:  "chain",
	Usage: "commands for chain",
	Subcommands: []cli.Command{
		cmdChainGet,
	},
}
