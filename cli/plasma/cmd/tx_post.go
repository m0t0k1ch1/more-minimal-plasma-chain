package cmd

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/client"
	"github.com/urfave/cli"
)

var CmdTxPost = cli.Command{
	Name:  "post",
	Usage: "post tx",
	Flags: []cli.Flag{
		hostFlag,
		txFlag,
	},
	Action: func(c *cli.Context) error {
		hostStr := c.String("host")
		txStr := c.String("tx")

		tx, err := decodeTx(txStr)
		if err != nil {
			return err
		}

		if _, err := client.New(hostStr).PostTx(context.Background(), tx); err != nil {
			return err
		}

		return nil
	},
}
