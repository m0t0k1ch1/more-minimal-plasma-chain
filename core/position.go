package core

import (
	"math/big"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

var (
	TxPositionOffset = big.NewInt(10000)
)

func TxPosition(blkNum, txIndex *big.Int) *big.Int {
	pos := new(big.Int).Mul(blkNum, big.NewInt(types.MaxBlockTxesNum+1))
	pos.Add(pos, txIndex)
	return pos
}

func TxOutPosition(blkNum, txIndex, oIndex *big.Int) *big.Int {
	pos := TxPosition(blkNum, txIndex)
	pos.Mul(pos, TxPositionOffset)
	pos.Add(pos, oIndex)
	return pos
}
