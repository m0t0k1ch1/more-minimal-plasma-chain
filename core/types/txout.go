package types

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
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

func (txOut *TxOut) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		txOut.OwnerAddress.Bytes(), txOut.Amount,
	})
}

func (txOut *TxOut) Spent() {
	txOut.IsSpent = true
}
