package main

import "github.com/urfave/cli"

var cmdRootChain = cli.Command{
	Name:  "rootchain",
	Usage: "commands for root chain",
	Subcommands: []cli.Command{
		cmdRootChainCommit,
	},
}
