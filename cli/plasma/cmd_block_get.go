package main

import (
	"context"

	"github.com/urfave/cli"
)

var cmdBlockGet = cli.Command{
	Name:  "get",
	Usage: "get block",
	Flags: []cli.Flag{
		apiFlag,
		hashFlag,
		encodedFlag,
	},
	Action: func(c *cli.Context) error {
		blkHash, err := getHash(c, hashFlag)
		if err != nil {
			return err
		}

		blk, err := newClient(c).GetBlock(
			context.Background(),
			blkHash,
		)
		if err != nil {
			return err
		}

		if getBool(c, encodedFlag) {
			return printlnEncodedBlock(blk)
		}
		return printlnJSON(blk)
	},
}
