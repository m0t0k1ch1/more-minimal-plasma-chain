package types

import (
	"github.com/ethereum/go-ethereum/common"
)

type TxOut struct {
	OwnerAddress common.Address `json:"owner"`
	Amount       uint64         `json:"amount"`
	IsSpent      bool           `json:"spent"`
}

func NewTxOut(ownerAddr common.Address, amount uint64) *TxOut {
	return &TxOut{
		OwnerAddress: ownerAddr,
		Amount:       amount,
		IsSpent:      false,
	}
}

func (txOut *TxOut) Spent() {
	txOut.IsSpent = true
}
