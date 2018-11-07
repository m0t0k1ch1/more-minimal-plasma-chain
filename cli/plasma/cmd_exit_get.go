package main

import (
	"github.com/urfave/cli"
)

var cmdExitGet = cli.Command{
	Name:  "get",
	Usage: "get exit",
	Flags: flags(
		posFlag,
	),
	Action: func(c *cli.Context) error {
		txOutPos, err := getPosition(c, posFlag)
		if err != nil {
			return err
		}

		rc, err := newRootChain()
		if err != nil {
			return err
		}

		exit, err := rc.PlasmaExits(txOutPos)
		if err != nil {
			return err
		}

		return printlnJSON(exit)
	},
}
