package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTxOutPosition(t *testing.T) {
	type input struct {
		blkNum   uint64
		txIndex  uint64
		outIndex uint64
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
				blkNum:   1,
				txIndex:  23456,
				outIndex: 7890,
			},
			output{
				pos: Position(1234567890),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			pos := NewTxOutPosition(in.blkNum, in.txIndex, in.outIndex)
			assert.Equal(t, out.pos, pos)
		})
	}
}

func TestParseTxOutPosition(t *testing.T) {
	type input struct {
		pos Position
	}
	type output struct {
		blkNum   uint64
		txIndex  uint64
		outIndex uint64
	}
	testCases := []struct {
		name string
		in   input
		out  output
	}{
		{
			"1234567890",
			input{
				Position(1234567890),
			},
			output{
				1,
				23456,
				7890,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			blkNum, txIndex, outIndex := ParseTxOutPosition(in.pos)
			assert.Equal(t, out.blkNum, blkNum)
			assert.Equal(t, out.txIndex, txIndex)
			assert.Equal(t, out.outIndex, outIndex)
		})
	}
}
