package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var cmdExitStart = cli.Command{
	Name:  "start",
	Usage: "start exit",
	Flags: flags(
		posFlag,
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		txOutPos, err := getPosition(c, posFlag)
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

		clnt := newClient()
		ctx := context.Background()

		blkNum, txIndex, _ := types.ParseTxOutPosition(txOutPos)
		txPos := types.NewTxPosition(blkNum, txIndex)

		// get tx
		tx, err := clnt.GetTx(ctx, txPos)
		if err != nil {
			return err
		}

		// get tx proof
		txProofBytes, err := clnt.GetTxProof(ctx, txPos)
		if err != nil {
			return err
		}

		// start exit
		txn, err := rc.StartExit(types.NewAccount(privKey), txOutPos, tx, txProofBytes)
		if err != nil {
			return err
		}

		return printlnJSON(txn)
	},
}
