package main

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdDepositMake = cli.Command{
	Name:  "make",
	Usage: "make deposit",
	Flags: flags(
		amountFlag,
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		amount, err := getUint64(c, amountFlag)
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

		rctx, err := rc.Deposit(types.NewAccount(privKey), amount)
		if err != nil {
			return err
		}

		return printlnJSON(rctx)
	},
}
