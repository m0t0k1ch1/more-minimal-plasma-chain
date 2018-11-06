package main

import (
	"context"

	"github.com/urfave/cli"
)

var cmdExitGet = cli.Command{
	Name:  "get",
	Usage: "get exit",
	Flags: []cli.Flag{
		rpcFlag, wsFlag, contractAddrFlag,
		apiFlag,
		txHashFlag,
		oIndexFlag,
	},
	Action: func(c *cli.Context) error {
		txHash, err := getHash(c, txHashFlag)
		if err != nil {
			return err
		}
		oIndex, err := getBigInt(c, oIndexFlag)
		if err != nil {
			return err
		}

		rc, err := newRootChain(c)
		if err != nil {
			return err
		}

		// get tx index
		blkNum, txIndex, err := newClient(c).GetTxIndex(context.Background(), txHash)
		if err != nil {
			return err
		}

		// get exit
		exit, err := rc.PlasmaExits(blkNum, txIndex, oIndex)
		if err != nil {
			return err
		}

		return printlnJSON(exit)
	},
}
