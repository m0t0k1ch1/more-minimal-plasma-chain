package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeToUtxoPosition(t *testing.T) {
	type input struct {
		blkNum  uint64
		txIndex uint64
		oIndex  uint64
	}
	type output struct {
		utxoPos uint64
	}
	testCases := []struct {
		name string
		in   input
		out  output
	}{
		{
			"1-1-1",
			input{
				1, 1, 1,
			},
			output{
				1000010001,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			utxoPos := EncodeToUtxoPosition(in.blkNum, in.txIndex, in.oIndex)
			assert.Equal(t, out.utxoPos, utxoPos)
		})
	}
}

func TestDecodeUtxoPosition(t *testing.T) {
	type input struct {
		utxoPos uint64
	}
	type output struct {
		blkNum  uint64
		txIndex uint64
		oIndex  uint64
	}
	testCases := []struct {
		name string
		in   input
		out  output
	}{
		{
			"1000010001",
			input{
				1000010001,
			},
			output{
				1, 1, 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			blkNum, txIndex, oIndex := DecodeUtxoPosition(in.utxoPos)
			assert.Equal(t, out.blkNum, blkNum)
			assert.Equal(t, out.txIndex, txIndex)
			assert.Equal(t, out.oIndex, oIndex)
		})
	}
}
