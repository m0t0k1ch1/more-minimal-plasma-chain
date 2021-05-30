package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdDepositMake = cli.Command{
	Name:  "make",
	Usage: "make deposit",
	Flags: flags(
		amountFlag,
		privKeyFlag,
		directFlag,
	),
	Action: func(c *cli.Context) error {
		amount, err := getBigInt(c, amountFlag)
		if err != nil {
			return err
		}
		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}
		isDirect := getBool(c, directFlag)

		account := types.NewAccount(privKey)

		if isDirect {
			// deposit to child chain directly
			blkNum, err := newClient().PostDeposit(context.Background(), account.Address(), amount)
			if err != nil {
				return err
			}

			return printlnJSON(map[string]uint64{
				"blknum": blkNum,
			})
		}

		rc, err := newRootChain()
		if err != nil {
			return err
		}

		// deposit to root chain
		rctx, err := rc.Deposit(account, amount)
		if err != nil {
			return err
		}

		return printlnJSON(rctx)
	},
}
