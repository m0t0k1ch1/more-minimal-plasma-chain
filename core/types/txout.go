package types

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	nullTxOut = NewTxOut(NullAddress, 0)
)

type TxOut struct {
	OwnerAddress common.Address `json:"owner"`
	Amount       uint64         `json:"amount"`
}

func NewTxOut(ownerAddr common.Address, amount uint64) *TxOut {
	return &TxOut{
		OwnerAddress: ownerAddr,
		Amount:       amount,
	}
}

func (txOut *TxOut) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		txOut.OwnerAddress.Bytes(), txOut.Amount,
	})
}
