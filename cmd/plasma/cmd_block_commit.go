package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var cmdBlockCommit = cli.Command{
	Name:  "commit",
	Usage: "commit block root to root chain",
	Flags: flags(
		numFlag,
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		blkNum, err := getUint64(c, numFlag)
		if err != nil {
			return err
		}
		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}

		rc, err := newRootChain()
		if err != nil {
			return err
		}

		// get block
		blk, err := newClient().GetBlock(context.Background(), blkNum)
		if err != nil {
			return err
		}

		// get block root hash
		blkRoot, err := blk.Root()
		if err != nil {
			return err
		}

		// commit block root hash
		rctx, err := rc.CommitPlasmaBlockRoot(types.NewAccount(privKey), utils.BytesToHash(blkRoot.Bytes()))
		if err != nil {
			return err
		}

		return printlnJSON(rctx)
	},
}
