package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var cmdBlockFix = cli.Command{
	Name:  "fix",
	Usage: "fix block",
	Flags: []cli.Flag{
		apiFlag,
	},
	Action: func(c *cli.Context) error {
		blkHash, err := newClient(c).PostBlock(context.Background())
		if err != nil {
			return err
		}

		return printlnJSON(map[string]string{
			"blkhash": utils.HashToHex(blkHash),
		})
	},
}
