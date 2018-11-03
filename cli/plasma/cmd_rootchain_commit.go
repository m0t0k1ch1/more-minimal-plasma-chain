package main

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdRootChainCommit = cli.Command{
	Name:  "commit",
	Usage: "commit block root",
	Flags: []cli.Flag{
		rpcFlag, wsFlag, contractFlag,
		privKeyFlag,
		rootFlag,
	},
	Action: func(c *cli.Context) error {
		rootHash, err := getHash(c, rootFlag)
		if err != nil {
			return err
		}

		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}
		a := types.NewAccount(privKey)

		rc, err := newRootChain(c)
		if err != nil {
			return err
		}

		txn, err := rc.CommitPlasmaBlockRoot(a, rootHash)
		if err != nil {
			return err
		}

		return printlnJSON(txn)
	},
}
