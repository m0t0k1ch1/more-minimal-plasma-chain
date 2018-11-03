package main

import (
	"context"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var cmdTxProof = cli.Command{
	Name:  "proof",
	Usage: "get tx proof",
	Flags: []cli.Flag{
		apiFlag,
		hashFlag,
	},
	Action: func(c *cli.Context) error {
		txHashBytes, err := getHexBytes(c, hashFlag)
		if err != nil {
			return err
		}

		proofBytes, err := newClient(c).GetTxProof(context.Background(), txHashBytes)
		if err != nil {
			return err
		}

		return printlnJSON(map[string]string{
			"proof": utils.EncodeToHex(proofBytes),
		})
	},
}
