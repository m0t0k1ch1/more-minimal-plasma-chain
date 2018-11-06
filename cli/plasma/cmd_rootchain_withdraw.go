package main

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdRootChainWithdraw = cli.Command{
	Name:  "withdraw",
	Usage: "withdraw",
	Flags: []cli.Flag{
		rpcFlag, wsFlag, contractAddrFlag,
		privKeyFlag,
	},
	Action: func(c *cli.Context) error {
		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}

		rc, err := newRootChain(c)
		if err != nil {
			return err
		}

		txn, err := rc.ProcessExits(types.NewAccount(privKey))
		if err != nil {
			return err
		}

		return printlnJSON(txn)
	},
}
