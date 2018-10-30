package types

import (
	"io"

	"github.com/ethereum/go-ethereum/rlp"
)

var (
	nullTxIn = NewTxIn(0, 0, 0)
)

type TxIn struct {
	BlockNumber uint64 `json:"blknum"`
	TxIndex     uint64 `json:"txindex"`
	OutputIndex uint64 `json:"oindex"`
}

func NewTxIn(blkNum, txIndex, oIndex uint64) *TxIn {
	return &TxIn{
		BlockNumber: blkNum,
		TxIndex:     txIndex,
		OutputIndex: oIndex,
	}
}

func (txIn *TxIn) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex,
	})
}
