package main

import (
	"context"

	"github.com/urfave/cli"
)

var cmdBlockGet = cli.Command{
	Name:  "get",
	Usage: "get block",
	Flags: flags(
		numFlag,
		encodedFlag,
	),
	Action: func(c *cli.Context) error {
		blkNum, err := getUint64(c, numFlag)
		if err != nil {
			return err
		}

		blk, err := newClient().GetBlock(context.Background(), blkNum)
		if err != nil {
			return err
		}

		if getBool(c, encodedFlag) {
			return printlnEncodedBlock(blk)
		}
		return printlnJSON(blk)
	},
}
