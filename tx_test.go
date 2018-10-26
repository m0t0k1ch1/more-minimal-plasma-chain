package main

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTx_Hash(t *testing.T) {
	testCases := []struct {
		name string
		tx   *Tx
		hex  string
	}{
		{
			"null tx",
			NullTx,
			"c758c57a2f76021ff85aa579dc03dc81ee6302c88eb43fd73190f1b036e5f0e6",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.tx.Hash()
			require.NoError(t, err)
			assert.Equal(t, hex.EncodeToString(b), tc.hex)
		})
	}
}
