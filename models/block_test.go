package models

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
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

func TestBlock_Sign(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	addr := crypto.PubkeyToAddress(privKey.PublicKey)

	blk := NewBlock(0)
	require.NoError(t, blk.Sign(privKey))

	signer, err := blk.Signer()
	require.NoError(t, err)

	assert.Equal(t, addr, signer)
}
