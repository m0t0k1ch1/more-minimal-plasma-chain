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
		blkNum, err := getBigInt(c, blkNumFlag)
		if err != nil {
			return err
		}
		txIndex, err := getBigInt(c, txIndexFlag)
		if err != nil {
			return err
		}
		oIndex, err := getBigInt(c, oIndexFlag)
		if err != nil {
			return err
		}

		ownerAddr, err := getAddress(c, ownerFlag)
		if err != nil {
			return err
		}
		amount, err := getBigInt(c, amountFlag)
		if err != nil {
			return err
		}

		tx := types.NewTx()
		tx.Inputs[0] = types.NewTxIn(blkNum, txIndex, oIndex)
		tx.Outputs[0] = types.NewTxOut(ownerAddr, amount)

		return printlnEncodedTx(tx)
	},
}
