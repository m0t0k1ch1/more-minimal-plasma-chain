package main

import (
	"github.com/urfave/cli"
)

var CmdBlock = cli.Command{
	Name:  "block",
	Usage: "commands for block",
	Subcommands: []cli.Command{
		CmdBlockFix,
		CmdBlockGet,
	},
}
