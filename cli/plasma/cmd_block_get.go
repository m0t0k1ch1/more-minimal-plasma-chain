package main

import (
	"context"

	"github.com/urfave/cli"
)

var cmdBlockGet = cli.Command{
	Name:  "get",
	Usage: "get block",
	Flags: []cli.Flag{
		hostFlag,
		portFlag,
		hashFlag,
	},
	Action: func(c *cli.Context) error {
		blkHashBytes, err := getHexBytes(c, hashFlag)
		if err != nil {
			return err
		}

		blk, err := newClient(c).GetBlock(context.Background(), blkHashBytes)
		if err != nil {
			return err
		}

		return printlnJSON(blk)
	},
}
