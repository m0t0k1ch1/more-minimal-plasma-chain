package cmd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var CmdTx = cli.Command{
	Name:  "tx",
	Usage: "commands for tx",
	Subcommands: []cli.Command{
		CmdTxCreate,
		CmdTxSign,
		CmdTxPost,
		CmdTxConfirm,
	},
}

func decodeTx(txStr string) (*types.Tx, error) {
	txBytes, err := hexutil.Decode(txStr)
	if err != nil {
		return nil, err
	}

	var tx types.Tx
	if err := rlp.DecodeBytes(txBytes, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}

func printlnTx(tx *types.Tx) error {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}

	fmt.Println(hexutil.Encode(txBytes))

	return nil
}
