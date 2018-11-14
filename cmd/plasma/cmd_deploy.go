package main

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

var cmdDeploy = cli.Command{
	Name:  "deploy",
	Usage: "deploy contract",
	Flags: flags(
		privKeyFlag,
	),
	Action: func(c *cli.Context) error {
		privKey, err := getPrivateKey(c, privKeyFlag)
		if err != nil {
			return err
		}

		abi, err := abi.JSON(strings.NewReader(core.RootChainABI))
		if err != nil {
			return err
		}

		bin, err := utils.HexToBytes(core.RootChainBin)
		if err != nil {
			return err
		}

		clnt, err := ethclient.Dial(conf.RootChain.RPC)
		if err != nil {
			return err
		}

		// deploy root chain contract
		addr, _, _, err := bind.DeployContract(types.NewAccount(privKey).TransactOpts(), abi, bin, clnt)
		if err != nil {
			return err
		}

		return printlnJSON(map[string]string{
			"address": utils.AddressToHex(addr),
		})
	},
}
