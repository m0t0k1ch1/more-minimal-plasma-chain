package types

import (
	"github.com/ethereum/go-ethereum/common"
)

type TxOutCore struct {
	OwnerAddress common.Address `json:"owner"`
	Amount       uint64         `json:"amount"`
}

type TxOut struct {
	*TxOutCore
	IsSpent bool `json:"spent"`
}

func NewTxOut(ownerAddr common.Address, amount uint64) *TxOut {
	return &TxOut{
		TxOutCore: &TxOutCore{
			OwnerAddress: ownerAddr,
			Amount:       amount,
		},
		IsSpent: false,
	}
}
