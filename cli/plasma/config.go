package main

import "github.com/m0t0k1ch1/more-minimal-plasma-chain/core"

type Config struct {
	RootChain  core.RootChainConfig `json:"rootchain"`
	ChildChain ChildChainConfig     `json:"childchain"`
}

type ChildChainConfig struct {
	API string `json:"api"`
}
