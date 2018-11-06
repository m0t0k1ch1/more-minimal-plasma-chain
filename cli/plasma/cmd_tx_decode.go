package main

import "github.com/urfave/cli"

var cmdTxDecode = cli.Command{
	Name:  "decode",
	Usage: "decode tx",
	Flags: []cli.Flag{
		txFlag,
	},
	Action: func(c *cli.Context) error {
		tx, err := getTx(c, txFlag)
		if err != nil {
			return err
		}

		return printlnJSON(tx)
	},
}
