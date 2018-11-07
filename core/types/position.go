package types

import (
	"math/big"
)

const (
	BlockPositionOffset = 100000 // must be greater than MaxBlockTxesNum
	TxPositionOffset    = 10000  // must be greater than TxElementsNum
)

var (
	NullPosition = NewPosition(nil)
)

type Position struct {
	*big.Int
}

func NewPosition(i *big.Int) Position {
	return Position{i}
}

func NewTxPosition(blkNum, txIndex *big.Int) Position {
	pos := new(big.Int).Mul(blkNum, big.NewInt(BlockPositionOffset))
	pos.Add(pos, txIndex)
	return NewPosition(pos)
}

func NewTxInPosition(blkNum, txIndex, inIndex *big.Int) Position {
	return newTxElementPosition(blkNum, txIndex, inIndex)
}

func NewTxOutPosition(blkNum, txIndex, outIndex *big.Int) Position {
	return newTxElementPosition(blkNum, txIndex, outIndex)
}

func newTxElementPosition(blkNum, txIndex, index *big.Int) Position {
	pos := NewTxPosition(blkNum, txIndex)
	pos.Mul(pos.Int, big.NewInt(TxPositionOffset))
	pos.Add(pos.Int, index)
	return pos
}

func ParseTxPosition(pos Position) (blkNum, txIndex *big.Int) {
	blkNum = new(big.Int).Div(pos.Int, big.NewInt(BlockPositionOffset))
	txIndex = new(big.Int).Mod(pos.Int, big.NewInt(BlockPositionOffset))
	return
}

func ParseTxInPosition(pos Position) (*big.Int, *big.Int, *big.Int) {
	return parseTxElementPosition(pos)
}

func ParseTxOutPosition(pos Position) (*big.Int, *big.Int, *big.Int) {
	return parseTxElementPosition(pos)
}

func parseTxElementPosition(pos Position) (blkNum, txIndex, index *big.Int) {
	txPos := new(big.Int).Div(pos.Int, big.NewInt(TxPositionOffset))
	blkNum, txIndex = ParseTxPosition(NewPosition(txPos))
	index = new(big.Int).Mod(pos.Int, txPos)
	return
}
