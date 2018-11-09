package main

import (
	"context"
	"fmt"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdTxCreate = cli.Command{
	Name:  "create",
	Usage: "create tx",
	Flags: flags(
		posFlag,
		toFlag,
		amountFlag,
	),
	Action: func(c *cli.Context) error {
		txOutPos, err := getPosition(c, posFlag)
		if err != nil {
			return err
		}
		toAddr, err := getAddress(c, toFlag)
		if err != nil {
			return err
		}
		amount, err := getUint64(c, amountFlag)
		if err != nil {
			return err
		}

		clnt := newClient()
		ctx := context.Background()

		blkNum, txIndex, outIndex := types.ParseTxOutPosition(txOutPos)
		txPos := types.NewTxPosition(blkNum, txIndex)

		// get input tx
		inTx, err := clnt.GetTx(ctx, txPos)
		if err != nil {
			return err
		}

		// get input UTXO
		inTxOut := inTx.GetOutput(outIndex)

		// validate amount
		if amount > inTxOut.Amount {
			return fmt.Errorf("invalid amount")
		}

		// calculate change amount
		changeAmount := inTxOut.Amount - amount

		// create tx
		tx := types.NewTx()
		if err := tx.SetInput(0, types.NewTxIn(blkNum, txIndex, outIndex)); err != nil {
			return err
		}
		if err := tx.SetOutput(0, types.NewTxOut(toAddr, amount)); err != nil {
			return err
		}
		if changeAmount > 0 {
			if err := tx.SetOutput(1, types.NewTxOut(inTxOut.OwnerAddress, changeAmount)); err != nil {
				return err
			}
		}

		return printlnEncodedTx(tx)
	},
}
