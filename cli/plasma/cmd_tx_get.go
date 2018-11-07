package main

import (
	"context"

	"github.com/urfave/cli"
)

var cmdTxGet = cli.Command{
	Name:  "get",
	Usage: "get tx",
	Flags: flags(
		txPosFlag,
		encodedFlag,
	),
	Action: func(c *cli.Context) error {
		txPos, err := getPosition(c, txPosFlag)
		if err != nil {
			return err
		}

		tx, err := newClient().GetTx(context.Background(), txPos)
		if err != nil {
			return err
		}

		if getBool(c, encodedFlag) {
			return printlnEncodedTx(tx)
		}
		return printlnJSON(tx)
	},
}
