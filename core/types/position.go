package types

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

const (
	BlockPositionOffset = 100000 // must be greater than MaxBlockTxesNum
	TxPositionOffset    = 10000  // must be greater than TxElementsNum
)

type Position uint64

func (pos Position) Uint64() uint64 {
	return uint64(pos)
}

func (pos Position) Bytes() []byte {
	return PositionToBytes(pos)
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

func PositionToBytes(pos Position) []byte {
	return utils.Uint64ToBytes(pos.Uint64())
}

func BytesToPosition(b []byte) (Position, error) {
	i, err := utils.BytesToUint64(b)
	if err != nil {
		return 0, err
	}
	return Position(i), nil
}

func StrToPosition(s string) (Position, error) {
	i, err := utils.StringToUint64(s)
	if err != nil {
		return 0, err
	}
	return Position(i), nil
}
