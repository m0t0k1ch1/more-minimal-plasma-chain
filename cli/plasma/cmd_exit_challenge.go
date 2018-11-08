package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdExitChallenge = cli.Command{
	Name:  "challenge",
	Usage: "challenge exit",
	Flags: flags(
		posFlag,
		vsPosFlag,
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		txOutPos, err := getPosition(c, posFlag)
		if err != nil {
			return err
		}
		spendingTxInPos, err := getPosition(c, vsPosFlag)
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

		spendingBlkNum, spendingTxIndex, spendingInIndex := types.ParseTxInPosition(spendingTxInPos)
		spendingTxPos := types.NewTxPosition(spendingBlkNum, spendingTxIndex)

		// get spending tx
		spendingTx, err := newClient().GetTx(context.Background(), spendingTxPos)
		if err != nil {
			return err
		}

		// challenge exit
		rctx, err := rc.ChallengeExit(types.NewAccount(privKey), txOutPos, spendingTx, spendingInIndex)
		if err != nil {
			return err
		}

		return printlnJSON(rctx)
	},
}
