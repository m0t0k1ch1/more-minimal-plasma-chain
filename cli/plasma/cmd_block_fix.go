package main

import (
	"context"
	"math/big"

	"github.com/urfave/cli"
)

var cmdBlockFix = cli.Command{
	Name:  "fix",
	Usage: "fix block",
	Flags: flags(),
	Action: func(c *cli.Context) error {
		blkNum, err := newClient().PostBlock(context.Background())
		if err != nil {
			return err
		}

		return printlnJSON(map[string]*big.Int{
			"blknum": blkNum,
		})
	},
}
