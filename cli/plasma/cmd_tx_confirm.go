package main

import (
	"context"
	"math/big"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var cmdTxConfirm = cli.Command{
	Name:  "confirm",
	Usage: "confirm tx",
	Flags: []cli.Flag{
		apiFlag,
		hashFlag,
		privKeyFlag,
	},
	Action: func(c *cli.Context) error {
		txHash, err := getHash(c, hashFlag)
		if err != nil {
			return err
		}

		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}

		ctx := context.Background()
		zero := big.NewInt(0)

		// get tx
		tx, err := newClient(c).GetTx(
			ctx,
			txHash,
		)
		if err != nil {
			return err
		}

		// confirm tx
		if err := tx.Confirm(zero, types.NewAccount(privKey)); err != nil {
			return err
		}

		// update confirmation signature
		if _, err := newClient(c).PutTx(
			ctx,
			txHash, zero, tx.Inputs[0].ConfirmationSignature,
		); err != nil {
			return err
		}

		return printlnJSON(map[string]string{
			"txhash": utils.HashToHex(txHash),
		})
	},
}
