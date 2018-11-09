package main

import (
	"context"

	"github.com/urfave/cli"
)

var cmdTxPost = cli.Command{
	Name:  "post",
	Usage: "post tx",
	Flags: flags(
		txFlag,
	),
	Action: func(c *cli.Context) error {
		tx, err := getTx(c, txFlag)
		if err != nil {
			return err
		}

		if err := newClient().PostTx(context.Background(), tx); err != nil {
			return err
		}

		return printlnJSON(nil)
	},
}
