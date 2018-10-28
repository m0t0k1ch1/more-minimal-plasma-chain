package types

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	nullTxOut = NewTxOut(nullAddress, 0)
)

type TxOut struct {
	Owner  common.Address `json:"owner"`
	Amount uint64         `json:"amount"`
}

func NewTxOut(owner common.Address, amount uint64) *TxOut {
	return &TxOut{
		Owner:  owner,
		Amount: amount,
	}
}

func (txOut *TxOut) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		txOut.Owner.Bytes(), txOut.Amount,
	})
}
