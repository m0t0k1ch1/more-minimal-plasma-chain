package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdTxOutRestore = cli.Command{
	Name:  "restore",
	Usage: "restore txout",
	Flags: flags(
		posFlag,
	),
	Action: func(c *cli.Context) error {
		txOutPos, err := getPosition(c, posFlag)
		if err != nil {
			return err
		}

		clnt := newClient()
		ctx := context.Background()

		// restore txout
		if err := clnt.PutTxOut(ctx, txOutPos, false); err != nil {
			return err
		}

		blkNum, txIndex, _ := types.ParseTxOutPosition(txOutPos)
		txPos := types.NewTxPosition(blkNum, txIndex)

		// get tx
		tx, err := clnt.GetTx(ctx, txPos)
		if err != nil {
			return err
		}

		return printlnJSON(tx)
	},
}
