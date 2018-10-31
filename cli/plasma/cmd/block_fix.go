package cmd

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/client"
	"github.com/urfave/cli"
)

var CmdBlockFix = cli.Command{
	Name:  "fix",
	Usage: "fix block",
	Flags: []cli.Flag{
		hostFlag,
	},
	Action: func(c *cli.Context) error {
		hostStr := c.String("host")

		blkNum, err := client.New(hostStr).PostBlock(context.Background())
		if err != nil {
			return err
		}

		return printlnJSON(map[string]interface{}{
			"blknum": blkNum,
		})
	},
}
