package cmd

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var CmdTxCreate = cli.Command{
	Name:  "create",
	Usage: "create tx",
	Flags: []cli.Flag{
		cli.Uint64Flag{Name: "blknum", Value: 0},
		cli.Uint64Flag{Name: "txindex", Value: 0},
		cli.Uint64Flag{Name: "oindex", Value: 0},
		cli.StringFlag{Name: "owner", Value: NullAddressStr},
		cli.Uint64Flag{Name: "amount", Value: 0},
	},
	Action: func(c *cli.Context) error {
		blkNum := c.Uint64("blknum")
		txIndex := c.Uint64("txindex")
		oIndex := c.Uint64("oindex")
		ownerAddrStr := c.String("owner")
		amount := c.Uint64("amount")

		ownerAddr, err := decodeAddress(ownerAddrStr)
		if err != nil {
			return err
		}

		tx := types.NewTx()
		tx.Inputs[0] = types.NewTxIn(blkNum, txIndex, oIndex)
		tx.Outputs[0] = types.NewTxOut(ownerAddr, amount)

		return printlnTx(tx)
	},
}
