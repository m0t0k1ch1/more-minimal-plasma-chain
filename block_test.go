package main

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlock_Hash(t *testing.T) {
	testCases := []struct {
		name string
		blk  *Block
		hex  string
	}{
		{
			"null block",
			NewBlock(0),
			"122d07b601c05953fe8229d17e5b5c0a66fbec3b9da839aea24afc18d86a6219",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.blk.Hash()
			require.NoError(t, err)
			assert.Equal(t, tc.hex, hex.EncodeToString(b))
		})
	}
}
