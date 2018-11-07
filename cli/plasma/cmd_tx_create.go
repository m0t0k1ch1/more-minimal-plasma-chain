package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdTxCreate = cli.Command{
	Name:  "create",
	Usage: "create tx",
	Flags: flags(
		txOutPosFlag,
		toAddrFlag,
		amountFlag,
	),
	Action: func(c *cli.Context) error {
		txOutPos, err := getPosition(c, txOutPosFlag)
		if err != nil {
			return err
		}
		toAddr, err := getAddress(c, toAddrFlag)
		if err != nil {
			return err
		}
		amount, err := getBigInt(c, amountFlag)
		if err != nil {
			return err
		}

		clnt := newClient()
		ctx := context.Background()

		blkNum, txIndex, oIndex := types.ParseTxOutPosition(txOutPos)
		txPos := types.NewTxPosition(blkNum, txIndex)

		// get input tx
		inTx, err := clnt.GetTx(ctx, txPos)
		if err != nil {
			return err
		}

		// get input UTXO
		inTxOut := inTx.GetOutput(oIndex)

		// validate amount
		if amount.Cmp(inTxOut.Amount) > 0 {
			return fmt.Errorf("invalid amount")
		}

		// calculate change amount
		changeAmount := new(big.Int).Sub(inTxOut.Amount, amount)

		// create tx
		tx := types.NewTx()
		tx.Inputs[0] = types.NewTxIn(blkNum, txIndex, oIndex)
		tx.Outputs[0] = types.NewTxOut(toAddr, amount)
		if changeAmount.Cmp(big.NewInt(0)) > 0 {
			tx.Outputs[1] = types.NewTxOut(inTxOut.OwnerAddress, changeAmount)
		}

		return printlnEncodedTx(tx)
	},
}
