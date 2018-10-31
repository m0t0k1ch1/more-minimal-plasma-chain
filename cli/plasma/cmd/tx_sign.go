package cmd

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var CmdTxSign = cli.Command{
	Name:  "sign",
	Usage: "sign tx",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "tx", Value: ""},
		cli.StringFlag{Name: "privkey", Value: ""},
	},
	Action: func(c *cli.Context) error {
		txStr := c.String("tx")
		privKeyStr := c.String("privkey")

		tx, err := decodeTx(txStr)
		if err != nil {
			return err
		}

		privKey, err := decodePrivateKey(privKeyStr)
		if err != nil {
			return err
		}
		signer := types.NewAccount(privKey)

		if err := tx.Sign(0, signer); err != nil {
			return err
		}

		return printlnTx(tx)
	},
}
