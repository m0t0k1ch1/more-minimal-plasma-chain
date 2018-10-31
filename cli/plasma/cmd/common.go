package cmd

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

const (
	NullAddressStr = "0x0000000000000000000000000000000000000000"
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
