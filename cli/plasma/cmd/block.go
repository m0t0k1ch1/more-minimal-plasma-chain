package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

var CmdBlock = cli.Command{
	Name:  "block",
	Usage: "commands for block",
	Subcommands: []cli.Command{
		CmdBlockFix,
		CmdBlockGet,
	},
}

func printlnBlockNumber(blkNum uint64) error {
	fmt.Println(map[string]interface{}{
		"blknum": blkNum,
	})

	return nil
}

func printlnBlock(blk *types.Block) error {
	blkBytes, err := json.Marshal(blk)
	if err != nil {
		return err
	}

	fmt.Println(string(blkBytes))

	return nil
}
