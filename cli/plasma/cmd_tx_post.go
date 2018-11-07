package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
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

		txPos, err := newClient().PostTx(context.Background(), tx)
		if err != nil {
			return err
		}

		return printlnJSON(map[string]types.Position{
			"pos": txPos,
		})
	},
}
