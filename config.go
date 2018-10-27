package main

import "github.com/ethereum/go-ethereum/common"

type Config struct {
	Port               int    `json:"port"`
	OperatorAddressHex string `json:"operator_address_hex"`
}

func (conf *Config) OperatorAddress() common.Address {
	return common.HexToAddress(conf.OperatorAddressHex)
}
