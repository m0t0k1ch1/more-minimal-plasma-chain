package main

const (
	BlockNumberOffset = 1000000000
	TxIndexOffset     = 10000
)

func EncodeToUtxoPosition(blkNum, txIndex, oIndex uint64) uint64 {
	return (blkNum * BlockNumberOffset) + (txIndex * TxIndexOffset) + oIndex
}

func DecodeUtxoPosition(utxoPos uint64) (uint64, uint64, uint64) {
	blkNum := utxoPos / BlockNumberOffset
	txIndex := (utxoPos % BlockNumberOffset) / TxIndexOffset
	oIndex := utxoPos - (blkNum * BlockNumberOffset) - (txIndex * TxIndexOffset)

	return blkNum, txIndex, oIndex
}
