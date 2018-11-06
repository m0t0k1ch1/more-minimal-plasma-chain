package core

import (
	"math/big"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

var (
	BlockPositionOffset = big.NewInt(types.MaxBlockTxesNum + 1)
	TxPositionOffset    = big.NewInt(10000)
)

func TxPosition(blkNum, txIndex *big.Int) *big.Int {
	pos := new(big.Int).Mul(blkNum, BlockPositionOffset)
	pos.Add(pos, txIndex)
	return pos
}

func TxOutPosition(blkNum, txIndex, oIndex *big.Int) *big.Int {
	pos := TxPosition(blkNum, txIndex)
	pos.Mul(pos, TxPositionOffset)
	pos.Add(pos, oIndex)
	return pos
}

func ParseTxPosition(pos *big.Int) (blkNum, txIndex *big.Int) {
	blkNum = new(big.Int).Div(pos, BlockPositionOffset)
	txIndex = new(big.Int).Mod(pos, BlockPositionOffset)
	return
}

func ParseTxOutPosition(pos *big.Int) (blkNum, txIndex, oIndex *big.Int) {
	txPos := new(big.Int).Div(pos, TxPositionOffset)
	blkNum, txIndex = ParseTxPosition(txPos)
	oIndex = new(big.Int).Mod(pos, txPos)
	return
}
