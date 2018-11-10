package main

import "github.com/urfave/cli"

var cmdAddress = cli.Command{
	Name:  "address",
	Usage: "commands for address",
	Subcommands: []cli.Command{
		cmdAddressUTXOs,
	},
}
