package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var cmdChainGet = cli.Command{
	Name:  "get",
	Usage: "get block hash",
	Flags: flags(
		blkNumFlag,
	),
	Action: func(c *cli.Context) error {
		blkNum, err := getBigInt(c, blkNumFlag)
		if err != nil {
			return err
		}

		blkHash, err := newClient(c).GetChain(context.Background(), blkNum)
		if err != nil {
			return err
		}

		return printlnJSON(map[string]string{
			"blkhash": utils.HashToHex(blkHash),
		})
	},
}
