package cmd

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/client"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
	"github.com/urfave/cli"
)

const (
	nullAddressStr = "0x0000000000000000000000000000000000000000"
)

var (
	hostFlag = cli.StringFlag{Name: "host", Value: "127.0.0.1"}
	portFlag = cli.Uint64Flag{Name: "port", Value: 1323}

	hashFlag    = cli.StringFlag{Name: "hash", Value: ""}
	privKeyFlag = cli.StringFlag{Name: "privkey", Value: ""}
	txFlag      = cli.StringFlag{Name: "tx", Value: ""}

	blkNumFlag  = cli.Uint64Flag{Name: "blknum", Value: 0}
	txIndexFlag = cli.Uint64Flag{Name: "txindex", Value: 0}
	oIndexFlag  = cli.Uint64Flag{Name: "oindex", Value: 0}
	ownerFlag   = cli.StringFlag{Name: "owner", Value: nullAddressStr}
	amountFlag  = cli.Uint64Flag{Name: "amount", Value: 0}
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
	privKeyBytes, err := hexutil.Decode(c.String(f.GetName()))
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
