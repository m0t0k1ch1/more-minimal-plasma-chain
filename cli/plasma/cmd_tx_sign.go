package main

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdTxSign = cli.Command{
	Name:  "sign",
	Usage: "sign tx",
	Flags: flags(
		txFlag,
		indexFlag,
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		tx, err := getTx(c, txFlag)
		if err != nil {
			return err
		}
		inIndex, err := getUint64(c, indexFlag)
		if err != nil {
			return err
		}
		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}

		if err := tx.Sign(inIndex, types.NewAccount(privKey)); err != nil {
			return err
		}

		return printlnEncodedTx(tx)
	},
}
