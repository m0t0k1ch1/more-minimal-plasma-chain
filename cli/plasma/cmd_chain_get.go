package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var cmdChainGet = cli.Command{
	Name:  "get",
	Usage: "get chain",
	Flags: []cli.Flag{
		hostFlag,
		portFlag,
		blkNumFlag,
	},
	Action: func(c *cli.Context) error {
		blkNum := getUint64(c, blkNumFlag)

		blkHashBytes, err := newClient(c).GetChain(context.Background(), blkNum)
		if err != nil {
			return err
		}

		return printlnJSON(map[string]interface{}{
			"blkhash": utils.EncodeToHex(blkHashBytes),
		})
	},
}
