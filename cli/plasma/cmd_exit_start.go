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
		txHashFlag,
		oIndexFlag,
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		txHash, err := getHash(c, txHashFlag)
		if err != nil {
			return err
		}
		oIndex, err := getBigInt(c, oIndexFlag)
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

		clnt := newClient(c)
		ctx := context.Background()

		// get tx
		tx, err := clnt.GetTx(ctx, txHash)
		if err != nil {
			return err
		}

		// get tx index
		blkNum, txIndex, err := clnt.GetTxIndex(ctx, txHash)
		if err != nil {
			return err
		}

		// get tx proof
		txProofBytes, err := clnt.GetTxProof(ctx, txHash)
		if err != nil {
			return err
		}

		// start exit
		txn, err := rc.StartExit(types.NewAccount(privKey), blkNum, txIndex, oIndex, tx, txProofBytes)
		if err != nil {
			return err
		}

		return printlnJSON(txn)
	},
}
