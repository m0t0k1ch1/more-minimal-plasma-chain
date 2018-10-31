package cmd

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/urfave/cli"
)

const (
	NullAddressStr = "0x0000000000000000000000000000000000000000"
)

var (
	hostFlag = cli.StringFlag{Name: "host", Value: "http://127.0.0.1:1323"}

	privKeyFlag = cli.StringFlag{Name: "privkey", Value: ""}
	txFlag      = cli.StringFlag{Name: "tx", Value: ""}
	blkNumFlag  = cli.Uint64Flag{Name: "blknum", Value: 0}
	txIndexFlag = cli.Uint64Flag{Name: "txindex", Value: 0}
	oIndexFlag  = cli.Uint64Flag{Name: "oindex", Value: 0}
	ownerFlag   = cli.StringFlag{Name: "owner", Value: NullAddressStr}
	amountFlag  = cli.Uint64Flag{Name: "amount", Value: 0}
)

func decodeAddress(addrStr string) (common.Address, error) {
	if !common.IsHexAddress(addrStr) {
		return types.NullAddress, fmt.Errorf("owner address is invalid")
	}

	return common.HexToAddress(addrStr), nil
}

func decodePrivateKey(privKeyStr string) (*ecdsa.PrivateKey, error) {
	privKeyBytes, err := hexutil.Decode(privKeyStr)
	if err != nil {
		return nil, err
	}

	return crypto.ToECDSA(privKeyBytes)
}
