package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
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
		txes   []*Tx
		blkNum uint64
	}
	type output struct {
		rootHex string
		err     error
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
				"0xe026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e83636516",
				nil,
			},
		},
		{
			"deposit block",
			input{
				[]*Tx{newTestDepositTx(t)},
				1,
			},
			output{
				"0xf88f3819a6a679a60f8d5070af717bdfb41a87ab9eceb631136273928fb30560",
				nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			blk, err := NewBlock(in.txes, in.blkNum)
			if out.err != nil {
				assert.EqualError(t, err, out.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, out.rootHex, hexutil.Encode(blk.Root()))
			}
		})
	}
}

func TestBlock_Hash(t *testing.T) {
	type output struct {
		hashHex string
		err     error
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
				"0x122d07b601c05953fe8229d17e5b5c0a66fbec3b9da839aea24afc18d86a6219",
				nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			blk, out := tc.blk, tc.out

			hashBytes, err := blk.Hash()
			if out.err != nil {
				assert.EqualError(t, err, out.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, out.hashHex, hexutil.Encode(hashBytes))
			}
		})
	}
}

func TestBlock_Sign(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	signer := NewAccount(privKey)
	blk := newTestNullBlock(t)

	// sign
	require.NoError(t, blk.Sign(signer))

	// verify
	signerAddr, err := blk.SignerAddress()
	require.NoError(t, err)
	assert.Equal(t, signer.Address(), signerAddr)
}
