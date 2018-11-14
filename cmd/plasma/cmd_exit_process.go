package main

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdExitProcess = cli.Command{
	Name:  "process",
	Usage: "process exits",
	Flags: flags(
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}

		rc, err := newRootChain()
		if err != nil {
			return err
		}

		rctx, err := rc.ProcessExits(types.NewAccount(privKey))
		if err != nil {
			return err
		}

		return printlnJSON(rctx)
	},
}
