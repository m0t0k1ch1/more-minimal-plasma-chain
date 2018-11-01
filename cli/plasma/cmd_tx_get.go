package main

import (
	"context"

	"github.com/urfave/cli"
)

var CmdTxGet = cli.Command{
	Name:  "get",
	Usage: "get tx",
	Flags: []cli.Flag{
		hostFlag,
		portFlag,
		hashFlag,
	},
	Action: func(c *cli.Context) error {
		txHashBytes, err := getHexBytes(c, hashFlag)
		if err != nil {
			return err
		}

		tx, err := newClient(c).GetTx(context.Background(), txHashBytes)
		if err != nil {
			return err
		}

		return printlnJSON(tx)
	},
}
