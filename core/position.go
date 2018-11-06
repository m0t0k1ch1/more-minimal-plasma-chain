package core

import "math/big"

func TxPosition(blkNum, txIndex *big.Int) *big.Int {
	pos := new(big.Int).Mul(blkNum, big.NewInt(100000))
	pos.Add(pos, txIndex)
	return pos
}

func TxOutPosition(blkNum, txIndex, oIndex *big.Int) *big.Int {
	pos := TxPosition(blkNum, txIndex)
	pos.Mul(pos, big.NewInt(10000))
	pos.Add(pos, oIndex)
	return pos
}
