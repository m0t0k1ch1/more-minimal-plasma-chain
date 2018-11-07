package types

import "math/big"

type TxInCore struct {
	BlockNumber *big.Int `json:"blknum"`
	TxIndex     *big.Int `json:"txindex"`
	OutputIndex *big.Int `json:"oindex"`
}

type TxIn struct {
	*TxInCore
	Signature             Signature `json:"sig"`
	ConfirmationSignature Signature `json:"confsig"`
}

func NewTxIn(blkNum, txIndex, outIndex *big.Int) *TxIn {
	return &TxIn{
		TxInCore: &TxInCore{
			BlockNumber: blkNum,
			TxIndex:     txIndex,
			OutputIndex: outIndex,
		},
		Signature:             NullSignature,
		ConfirmationSignature: NullSignature,
	}
}

func (txIn *TxIn) IsNull() bool {
	zero := big.NewInt(0)
	return txIn.BlockNumber.Cmp(zero) == 0 &&
		txIn.TxIndex.Cmp(zero) == 0 &&
		txIn.OutputIndex.Cmp(zero) == 0
}
