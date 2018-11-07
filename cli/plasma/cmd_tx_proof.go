package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var cmdTxProof = cli.Command{
	Name:  "proof",
	Usage: "get tx proof",
	Flags: flags(
		posFlag,
	),
	Action: func(c *cli.Context) error {
		txPos, err := getPosition(c, posFlag)
		if err != nil {
			return err
		}

		txProofBytes, err := newClient().GetTxProof(context.Background(), txPos)
		if err != nil {
			return err
		}

		return printlnJSON(map[string]string{
			"proof": utils.EncodeToHex(txProofBytes),
		})
	},
}
