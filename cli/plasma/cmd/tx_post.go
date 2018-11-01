package cmd

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var CmdTxPost = cli.Command{
	Name:  "post",
	Usage: "post tx",
	Flags: []cli.Flag{
		hostFlag,
		portFlag,
		txFlag,
	},
	Action: func(c *cli.Context) error {
		tx, err := getTx(c, txFlag)
		if err != nil {
			return err
		}

		txHashBytes, err := newClient(c).PostTx(context.Background(), tx)
		if err != nil {
			return err
		}

		return printlnJSON(map[string]interface{}{
			"txhash": utils.EncodeToHex(txHashBytes),
		})
	},
}
