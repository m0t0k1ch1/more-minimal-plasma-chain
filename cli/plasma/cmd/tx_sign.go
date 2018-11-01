package cmd

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var CmdTxSign = cli.Command{
	Name:  "sign",
	Usage: "sign tx",
	Flags: []cli.Flag{
		txFlag,
		privKeyFlag,
	},
	Action: func(c *cli.Context) error {
		tx, err := getTx(c, txFlag)
		if err != nil {
			return err
		}

		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}

		if err := tx.Sign(0, types.NewAccount(privKey)); err != nil {
			return err
		}

		return printlnEncodedTx(tx)
	},
}
