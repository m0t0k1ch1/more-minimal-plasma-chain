package core

import "math/big"

var (
	BlockPositionOffset = big.NewInt(1000000000)
	TxPositionOffset    = big.NewInt(10000)
)

func BlockPosition(blkNum *big.Int) *big.Int {
	return new(big.Int).Mul(blkNum, BlockPositionOffset)
}

func TxPosition(blkNum, txIndex *big.Int) *big.Int {
	pos := BlockPosition(blkNum)
	pos.Add(pos, new(big.Int).Mul(txIndex, TxPositionOffset))
	return pos
}

func TxOutPosition(blkNum, txIndex, oIndex *big.Int) *big.Int {
	pos := TxPosition(blkNum, txIndex)
	pos.Add(pos, oIndex)
	return pos
}
