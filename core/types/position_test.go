package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxElementPosition(t *testing.T) {
	type input struct {
		blkNum  *big.Int
		txIndex *big.Int
		oIndex  *big.Int
	}
	type output struct {
		pos Position
	}
	testCases := []struct {
		name string
		in   input
		out  output
	}{
		{
			"1234567890",
			input{
				blkNum:  big.NewInt(1),
				txIndex: big.NewInt(23456),
				oIndex:  big.NewInt(7890),
			},
			output{
				pos: Position{big.NewInt(1234567890)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			pos := TxElementPosition(in.blkNum, in.txIndex, in.oIndex)
			assert.Equal(t, out.pos, pos)
		})
	}
}

func TestParseTxElementPosition(t *testing.T) {
	type input struct {
		pos Position
	}
	type output struct {
		blkNum  *big.Int
		txIndex *big.Int
		oIndex  *big.Int
	}
	testCases := []struct {
		name string
		in   input
		out  output
	}{
		{
			"1234567890",
			input{
				Position{big.NewInt(1234567890)},
			},
			output{
				big.NewInt(1),
				big.NewInt(23456),
				big.NewInt(7890),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			blkNum, txIndex, oIndex := ParseTxElementPosition(in.pos)
			assert.Equal(t, out.blkNum, blkNum)
			assert.Equal(t, out.txIndex, txIndex)
			assert.Equal(t, out.oIndex, oIndex)
		})
	}
}
