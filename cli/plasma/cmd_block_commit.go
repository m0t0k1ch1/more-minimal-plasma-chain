package main

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdBlockCommit = cli.Command{
	Name:  "commit",
	Usage: "commit block root to root chain",
	Flags: flags(
		blkRootHashFlag,
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		blkRootHash, err := getHash(c, blkRootHashFlag)
		if err != nil {
			return err
		}
		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}

		rc, err := newRootChain(c)
		if err != nil {
			return err
		}

		txn, err := rc.CommitPlasmaBlockRoot(types.NewAccount(privKey), blkRootHash)
		if err != nil {
			return err
		}

		return printlnJSON(txn)
	},
}
