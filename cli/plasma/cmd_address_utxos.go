package main

import (
	"context"

	"github.com/urfave/cli"
)

var cmdAddressUTXOs = cli.Command{
	Name:  "utxos",
	Usage: "get utxos in address",
	Flags: flags(
		addressFlag,
	),
	Action: func(c *cli.Context) error {
		addr, err := getAddress(c, addressFlag)
		if err != nil {
			return err
		}

		utxoPoses, err := newClient().GetAddressUTXOs(context.Background(), addr)
		if err != nil {
			return err
		}

		return printlnJSON(utxoPoses)
	},
}
