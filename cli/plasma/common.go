package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/app"
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

func newRootChain(c *cli.Context) (*app.RootChain, error) {
	addr, err := getAddress(c, addrFlag)
	if err != nil {
		return nil, err
	}

	return app.NewRootChain(&app.RootChainConfig{
		RPC:     c.String(rpcFlag.GetName()),
		WS:      c.String(wsFlag.GetName()),
		Address: utils.EncodeToHex(addr.Bytes()),
	})
}

func getBool(c *cli.Context, f cli.Flag) bool {
	return c.Bool(f.GetName())
}

func getUint64(c *cli.Context, f cli.Flag) uint64 {
	return c.Uint64(f.GetName())
}

func getHexBytes(c *cli.Context, f cli.Flag) ([]byte, error) {
	return utils.DecodeHex(c.String(f.GetName()))
}

func getAddress(c *cli.Context, f cli.Flag) (common.Address, error) {
	addrStr := c.String(f.GetName())

	if !common.IsHexAddress(addrStr) {
		return types.NullAddress, fmt.Errorf("invalid address hex")
	}

	return common.HexToAddress(addrStr), nil
}

func getHash(c *cli.Context, f cli.Flag) (common.Hash, error) {
	hashStr := c.String(f.GetName())

	hashBytes, err := utils.DecodeHex(hashStr)
	if err != nil {
		return types.NullHash, err
	}
	if len(hashBytes) != common.HashLength {
		return types.NullHash, fmt.Errorf("invalid hash hex")
	}

	return common.BytesToHash(hashBytes), nil
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

func printlnEncodedBlock(blk *types.Block) error {
	blkBytes, err := rlp.EncodeToBytes(blk)
	if err != nil {
		return err
	}

	return printlnJSON(map[string]string{
		"blk": utils.EncodeToHex(blkBytes),
	})
}

func printlnEncodedTx(tx *types.Tx) error {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}

	return printlnJSON(map[string]string{
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
