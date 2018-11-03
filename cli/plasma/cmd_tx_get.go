package main

import (
	"context"

	"github.com/urfave/cli"
)

var cmdTxGet = cli.Command{
	Name:  "get",
	Usage: "get tx",
	Flags: []cli.Flag{
		apiFlag,
		hashFlag,
		encodedFlag,
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

		if getBool(c, encodedFlag) {
			return printlnEncodedTx(tx)
		}
		return printlnJSON(tx)
	},
}
