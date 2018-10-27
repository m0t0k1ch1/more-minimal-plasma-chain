package models

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestNullBlock(t *testing.T) *Block {
	blk, err := NewBlock(nil, 0)
	require.NoError(t, err)
	return blk
}

func TestNewBlock(t *testing.T) {
	type input struct {
		txes []*Tx
		num  uint64
	}
	type output struct {
		root string
		err  error
	}
	testCases := []struct {
		name string
		in   input
		out  output
	}{
		{
			"null block",
			input{
				nil,
				0,
			},
			output{
				"e026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e83636516",
				nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			blk, err := NewBlock(in.txes, in.num)
			if out.err != nil {
				assert.EqualError(t, err, out.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, out.root, common.Bytes2Hex(blk.Root()))
			}

		})
	}
}

func TestBlock_Hash(t *testing.T) {
	type output struct {
		hash string
		err  error
	}
	testCases := []struct {
		name string
		blk  *Block
		out  output
	}{
		{
			"null block",
			newTestNullBlock(t),
			output{
				"122d07b601c05953fe8229d17e5b5c0a66fbec3b9da839aea24afc18d86a6219",
				nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			blk, out := tc.blk, tc.out

			b, err := blk.Hash()
			if out.err != nil {
				assert.EqualError(t, err, out.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, out.hash, common.Bytes2Hex(b))
			}
		})
	}
}

func TestBlock_Sign(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	blk := newTestNullBlock(t)

	// sign
	require.NoError(t, blk.Sign(privKey))

	// verify
	signer, err := blk.Signer()
	require.NoError(t, err)
	assert.Equal(t, crypto.PubkeyToAddress(privKey.PublicKey), signer)
}
