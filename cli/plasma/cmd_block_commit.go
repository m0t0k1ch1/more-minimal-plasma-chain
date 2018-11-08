package main

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdBlockCommit = cli.Command{
	Name:  "commit",
	Usage: "commit block root to root chain",
	Flags: flags(
		rootFlag,
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		blkRootHash, err := getHash(c, rootFlag)
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

		rctx, err := rc.CommitPlasmaBlockRoot(types.NewAccount(privKey), blkRootHash)
		if err != nil {
			return err
		}

		return printlnJSON(rctx)
	},
}
