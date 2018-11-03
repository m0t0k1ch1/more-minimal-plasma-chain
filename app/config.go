package app

import "github.com/m0t0k1ch1/more-minimal-plasma-chain/core"

type Config struct {
	Port      int                   `json:"port"`
	Operator  *OperatorConfig       `json:"operator"`
	RootChain *core.RootChainConfig `json:"rootchain"`
}

type OperatorConfig struct {
	PrivateKey string `json:"privkey"`
}
