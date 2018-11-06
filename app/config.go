package app

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
)

type Config struct {
	Port      int                  `json:"port"`
	Operator  OperatorConfig       `json:"operator"`
	RootChain core.RootChainConfig `json:"rootchain"`
}

type OperatorConfig struct {
	PrivateKeyStr string `json:"privkey"`
}

func (conf OperatorConfig) PrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(conf.PrivateKeyStr)
}
