package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var cmdTxConfirm = cli.Command{
	Name:  "confirm",
	Usage: "confirm tx",
	Flags: flags(
		txHashFlag,
		iIndexFlag,
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		txHash, err := getHash(c, txHashFlag)
		if err != nil {
			return err
		}
		iIndex, err := getBigInt(c, iIndexFlag)
		if err != nil {
			return err
		}
		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}

		clnt := newClient()
		ctx := context.Background()

		// get tx
		tx, err := clnt.GetTx(ctx, txHash)
		if err != nil {
			return err
		}

		// confirm tx
		if err := tx.Confirm(iIndex, types.NewAccount(privKey)); err != nil {
			return err
		}

		// update confirmation signature
		if _, err := clnt.PutTx(ctx, txHash, iIndex, tx.GetInput(iIndex).ConfirmationSignature); err != nil {
			return err
		}

		return printlnJSON(map[string]string{
			"txhash": utils.HashToHex(txHash),
		})
	},
}
