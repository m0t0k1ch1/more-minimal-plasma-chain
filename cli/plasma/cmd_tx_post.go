package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var cmdTxPost = cli.Command{
	Name:  "post",
	Usage: "post tx",
	Flags: []cli.Flag{
		apiFlag,
		txFlag,
	},
	Action: func(c *cli.Context) error {
		tx, err := getTx(c, txFlag)
		if err != nil {
			return err
		}

		txHash, err := newClient(c).PostTx(
			context.Background(),
			tx,
		)
		if err != nil {
			return err
		}

		return printlnJSON(map[string]string{
			"txhash": utils.HashToHex(txHash),
		})
	},
}
