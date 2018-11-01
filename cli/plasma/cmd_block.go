package main

import (
	"github.com/urfave/cli"
)

var cmdBlock = cli.Command{
	Name:  "block",
	Usage: "commands for block",
	Subcommands: []cli.Command{
		cmdBlockFix,
		cmdBlockGet,
	},
}
