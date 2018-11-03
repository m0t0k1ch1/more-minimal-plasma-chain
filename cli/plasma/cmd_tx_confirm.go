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
	Flags: []cli.Flag{
		apiFlag,
		hashFlag,
		privKeyFlag,
	},
	Action: func(c *cli.Context) error {
		txHashBytes, err := getHexBytes(c, hashFlag)
		if err != nil {
			return err
		}

		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}

		tx, err := newClient(c).GetTx(context.Background(), txHashBytes)
		if err != nil {
			return err
		}

		if err := tx.Confirm(0, types.NewAccount(privKey)); err != nil {
			return err
		}

		if _, err := newClient(c).PutTx(context.Background(), txHashBytes, 0, tx.Inputs[0].ConfirmationSignature); err != nil {
			return err
		}

		return printlnJSON(map[string]string{
			"txhash": utils.EncodeToHex(txHashBytes),
		})
	},
}
