package main

import (
	"context"

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

		return printlnJSON(map[string]uint64{
			"blknum": blkNum,
		})
	},
}
