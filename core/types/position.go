package types

const (
	BlockPositionOffset = 100000 // must be greater than MaxBlockTxesNum
	TxPositionOffset    = 10000  // must be greater than TxElementsNum
)

type Position uint64

func (pos Position) Uint64() uint64 {
	return uint64(pos)
}

func NewTxPosition(blkNum, txIndex uint64) Position {
	return Position(blkNum*BlockPositionOffset + txIndex)
}

func NewTxInPosition(blkNum, txIndex, inIndex uint64) Position {
	return newTxElementPosition(blkNum, txIndex, inIndex)
}

func NewTxOutPosition(blkNum, txIndex, outIndex uint64) Position {
	return newTxElementPosition(blkNum, txIndex, outIndex)
}

func newTxElementPosition(blkNum, txIndex, elemIndex uint64) Position {
	return NewTxPosition(blkNum, txIndex)*TxPositionOffset + Position(elemIndex)
}

func ParseTxPosition(pos Position) (blkNum, txIndex uint64) {
	blkNum = pos.Uint64() / BlockPositionOffset
	txIndex = pos.Uint64() % BlockPositionOffset
	return
}

func ParseTxInPosition(pos Position) (uint64, uint64, uint64) {
	return parseTxElementPosition(pos)
}

func ParseTxOutPosition(pos Position) (uint64, uint64, uint64) {
	return parseTxElementPosition(pos)
}

func parseTxElementPosition(pos Position) (blkNum, txIndex, elemIndex uint64) {
	txPos := pos / Position(TxPositionOffset)
	blkNum, txIndex = ParseTxPosition(txPos)
	elemIndex = pos.Uint64() % txPos.Uint64()
	return
}
