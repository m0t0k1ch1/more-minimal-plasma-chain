package cmd

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/client"
	"github.com/urfave/cli"
)

var CmdBlockGet = cli.Command{
	Name:  "get",
	Usage: "get block",
	Flags: []cli.Flag{
		hostFlag,
		blkNumFlag,
	},
	Action: func(c *cli.Context) error {
		hostStr := c.String("host")
		blkNum := c.Uint64("blknum")

		blk, err := client.New(hostStr).GetBlock(context.Background(), blkNum)
		if err != nil {
			return err
		}

		return printlnBlock(blk)
	},
}
