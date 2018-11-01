package cmd

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var CmdBlockFix = cli.Command{
	Name:  "fix",
	Usage: "fix block",
	Flags: []cli.Flag{
		hostFlag,
		portFlag,
	},
	Action: func(c *cli.Context) error {
		blkHashBytes, err := newClient(c).PostBlock(context.Background())
		if err != nil {
			return err
		}

		return printlnJSON(map[string]interface{}{
			"blkhash": utils.EncodeToHex(blkHashBytes),
		})
	},
}
