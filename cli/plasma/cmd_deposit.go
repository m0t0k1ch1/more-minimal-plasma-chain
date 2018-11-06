package main

import "github.com/urfave/cli"

var cmdDeposit = cli.Command{
	Name:  "deposit",
	Usage: "commands for deposit",
	Subcommands: []cli.Command{
		cmdDepositMake,
	},
}
