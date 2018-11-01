package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/client"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

func newClient(c *cli.Context) *client.Client {
	return client.New(fmt.Sprintf(
		"http://%s:%d",
		c.String(hostFlag.GetName()), c.Uint64(portFlag.GetName()),
	))
}

func getUint64(c *cli.Context, f cli.Flag) uint64 {
	return c.Uint64(f.GetName())
}

func getHexBytes(c *cli.Context, f cli.Flag) ([]byte, error) {
	return utils.DecodeHex(c.String(f.GetName()))
}

func getAddress(c *cli.Context, f cli.Flag) (common.Address, error) {
	if !common.IsHexAddress(c.String(f.GetName())) {
		return types.NullAddress, fmt.Errorf("invalid address hex")
	}

	return common.HexToAddress(c.String(f.GetName())), nil
}

func getPrivateKey(c *cli.Context, f cli.Flag) (*ecdsa.PrivateKey, error) {
	privKeyBytes, err := utils.DecodeHex(c.String(f.GetName()))
	if err != nil {
		return nil, err
	}

	return crypto.ToECDSA(privKeyBytes)
}

func getTx(c *cli.Context, f cli.Flag) (*types.Tx, error) {
	txBytes, err := utils.DecodeHex(c.String(f.GetName()))
	if err != nil {
		return nil, err
	}

	var tx types.Tx
	if err := rlp.DecodeBytes(txBytes, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}

func printlnEncodedTx(tx *types.Tx) error {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}

	return printlnJSON(map[string]interface{}{
		"tx": utils.EncodeToHex(txBytes),
	})
}

func printlnJSON(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
