package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdTxPost = cli.Command{
	Name:  "post",
	Usage: "post tx",
	Flags: flags(
		posFlag,
		addressFlag,
		amountFlag,
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		txOutPos, err := getPosition(c, posFlag)
		if err != nil {
			return err
		}
		addr, err := getAddress(c, addressFlag)
		if err != nil {
			return err
		}
		amount, err := getBigInt(c, amountFlag)
		if err != nil {
			return err
		}
		privKey, err := getPrivateKey(c, privKeyFlag)
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
		if amount.Cmp(inTxOut.Amount) > 0 {
			return fmt.Errorf("invalid amount")
		}

		// calculate change amount
		changeAmount := new(big.Int).Sub(inTxOut.Amount, amount)

		// create tx
		tx := types.NewTx()
		if err := tx.SetInput(0, types.NewTxIn(blkNum, txIndex, outIndex)); err != nil {
			return err
		}
		if err := tx.SetOutput(0, types.NewTxOut(addr, amount)); err != nil {
			return err
		}
		if changeAmount.Sign() > 0 {
			if err := tx.SetOutput(1, types.NewTxOut(inTxOut.OwnerAddress, changeAmount)); err != nil {
				return err
			}
		}

		// sign tx
		if err := tx.Sign(0, types.NewAccount(privKey)); err != nil {
			return err
		}

		// post tx
		if err := clnt.PostTx(ctx, tx); err != nil {
			return err
		}

		return printlnJSON(tx)
	},
}
