package types

import (
	"math/big"
)

var (
	BlockPositionOffset = big.NewInt(MaxBlockTxesNum + 1)
	TxPositionOffset    = big.NewInt(10000)
)

type Position struct {
	*big.Int
}

func TxPosition(blkNum, txIndex *big.Int) Position {
	pos := new(big.Int).Mul(blkNum, BlockPositionOffset)
	pos.Add(pos, txIndex)
	return Position{pos}
}

func TxElementPosition(blkNum, txIndex, index *big.Int) Position {
	pos := TxPosition(blkNum, txIndex)
	pos.Mul(pos.Int, TxPositionOffset)
	pos.Add(pos.Int, index)
	return pos
}

func ParseTxPosition(pos Position) (blkNum, txIndex *big.Int) {
	blkNum = new(big.Int).Div(pos.Int, BlockPositionOffset)
	txIndex = new(big.Int).Mod(pos.Int, BlockPositionOffset)
	return
}

func ParseTxElementPosition(pos Position) (blkNum, txIndex, index *big.Int) {
	txPos := new(big.Int).Div(pos.Int, TxPositionOffset)
	blkNum, txIndex = ParseTxPosition(Position{txPos})
	index = new(big.Int).Mod(pos.Int, txPos)
	return
}
