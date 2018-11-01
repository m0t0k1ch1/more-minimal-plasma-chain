package main

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdTxCreate = cli.Command{
	Name:  "create",
	Usage: "create tx",
	Flags: []cli.Flag{
		blkNumFlag,
		txIndexFlag,
		oIndexFlag,
		ownerFlag,
		amountFlag,
	},
	Action: func(c *cli.Context) error {
		blkNum := getUint64(c, blkNumFlag)
		txIndex := getUint64(c, txIndexFlag)
		oIndex := getUint64(c, oIndexFlag)
		amount := getUint64(c, amountFlag)

		ownerAddr, err := getAddress(c, ownerFlag)
		if err != nil {
			return err
		}

		tx := types.NewTx()
		tx.Inputs[0] = types.NewTxIn(blkNum, txIndex, oIndex)
		tx.Outputs[0] = types.NewTxOut(ownerAddr, amount)

		return printlnEncodedTx(tx)
	},
}
