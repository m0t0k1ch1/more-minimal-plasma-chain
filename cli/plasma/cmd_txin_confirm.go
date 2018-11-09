package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdTxInConfirm = cli.Command{
	Name:  "confirm",
	Usage: "confirm tx input",
	Flags: flags(
		posFlag,
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		txInPos, err := getPosition(c, posFlag)
		if err != nil {
			return err
		}
		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}

		clnt := newClient()
		ctx := context.Background()

		blkNum, txIndex, inIndex := types.ParseTxInPosition(txInPos)
		txPos := types.NewTxPosition(blkNum, txIndex)

		// get tx
		tx, err := clnt.GetTx(ctx, txPos)
		if err != nil {
			return err
		}

		// confirm tx
		if err := tx.Confirm(inIndex, types.NewAccount(privKey)); err != nil {
			return err
		}

		// update confirmation signature
		if err := clnt.PutTxIn(ctx, txInPos, tx.GetInput(inIndex).ConfirmationSignature); err != nil {
			return err
		}

		return printlnJSON(map[string]types.Position{
			"pos": txInPos,
		})
	},
}
