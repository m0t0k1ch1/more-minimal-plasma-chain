package cmd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var CmdTxConfirm = cli.Command{
	Name:  "confirm",
	Usage: "confirm tx",
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

		if err := tx.Confirm(0, signer); err != nil {
			return err
		}

		fmt.Println(hexutil.Encode(tx.Inputs[0].ConfirmationSignature.Bytes()))

		return nil
	},
}
