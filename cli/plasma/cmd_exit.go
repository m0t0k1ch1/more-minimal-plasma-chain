package main

import (
	"github.com/urfave/cli"
)

var cmdExit = cli.Command{
	Name:  "exit",
	Usage: "commands for exit",
	Subcommands: []cli.Command{
		cmdExitGet,
		cmdExitProcess,
		cmdExitStart,
	},
}
