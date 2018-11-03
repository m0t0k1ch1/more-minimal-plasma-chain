package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TxOutCore struct {
	OwnerAddress common.Address `json:"owner"`
	Amount       *big.Int       `json:"amount"`
}

type TxOut struct {
	*TxOutCore
	IsSpent bool `json:"spent"`
}

func NewTxOut(ownerAddr common.Address, amount *big.Int) *TxOut {
	return &TxOut{
		TxOutCore: &TxOutCore{
			OwnerAddress: ownerAddr,
			Amount:       amount,
		},
		IsSpent: false,
	}
}
