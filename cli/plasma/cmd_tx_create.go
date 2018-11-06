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
	Flags: []cli.Flag{
		apiFlag,
		txHashFlag,
		oIndexFlag,
		toAddrFlag,
		amountFlag,
	},
	Action: func(c *cli.Context) error {
		txHash, err := getHash(c, txHashFlag)
		if err != nil {
			return err
		}
		oIndex, err := getBigInt(c, oIndexFlag)
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

		clnt := newClient(c)
		ctx := context.Background()

		// get input tx
		inTx, err := clnt.GetTx(
			ctx,
			txHash,
		)
		if err != nil {
			return err
		}

		// get input tx index
		blkNum, txIndex, err := clnt.GetTxIndex(
			ctx,
			txHash,
		)
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
