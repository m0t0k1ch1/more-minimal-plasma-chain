package types

type TxInCore struct {
	BlockNumber uint64 `json:"blknum"`
	TxIndex     uint64 `json:"txindex"`
	OutputIndex uint64 `json:"oindex"`
}

type TxIn struct {
	*TxInCore
	Signature             Signature `json:"sig"`
	ConfirmationSignature Signature `json:"confsig"`
}

func NewTxIn(blkNum, txIndex, oIndex uint64) *TxIn {
	return &TxIn{
		TxInCore: &TxInCore{
			BlockNumber: blkNum,
			TxIndex:     txIndex,
			OutputIndex: oIndex,
		},
		Signature:             NullSignature,
		ConfirmationSignature: NullSignature,
	}
}

func (txIn *TxIn) IsNull() bool {
	return txIn.BlockNumber == 0 && txIn.TxIndex == 0 && txIn.OutputIndex == 0
}
