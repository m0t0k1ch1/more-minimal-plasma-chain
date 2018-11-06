package main

import (
	"context"
	"math/big"

	"github.com/urfave/cli"
)

var cmdTxIndex = cli.Command{
	Name:  "index",
	Usage: "get tx index",
	Flags: flags(
		txHashFlag,
	),
	Action: func(c *cli.Context) error {
		txHash, err := getHash(c, txHashFlag)
		if err != nil {
			return err
		}

		blkNum, txIndex, err := newClient().GetTxIndex(context.Background(), txHash)
		if err != nil {
			return err
		}

		return printlnJSON(map[string]*big.Int{
			"blknum":  blkNum,
			"txindex": txIndex,
		})
	},
}
