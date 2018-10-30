package types

import (
	"io"

	"github.com/ethereum/go-ethereum/rlp"
)

type TxIn struct {
	BlockNumber           uint64    `json:"blknum"`
	TxIndex               uint64    `json:"txindex"`
	OutputIndex           uint64    `json:"oindex"`
	Signature             Signature `json:"sig"`
	ConfirmationSignature Signature `json:"confsig"`
}

func NewTxIn(blkNum, txIndex, oIndex uint64) *TxIn {
	return &TxIn{
		BlockNumber:           blkNum,
		TxIndex:               txIndex,
		OutputIndex:           oIndex,
		Signature:             NullSignature,
		ConfirmationSignature: NullSignature,
	}
}

func (txIn *TxIn) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex,
	})
}
