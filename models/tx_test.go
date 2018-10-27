package models

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testAddress = common.HexToAddress("0x82a978b3f5962a5b0957d9ee9eef472ee55b42f1")
)

func newTestDepositTx() *Tx {
	tx := NewTx()
	tx.Outputs[0] = NewTxOut(testAddress, 1)
	return tx
}

func TestTx_Hash(t *testing.T) {
	testCases := []struct {
		name string
		tx   *Tx
		hex  string
	}{
		{
			"null tx",
			NewTx(),
			"c758c57a2f76021ff85aa579dc03dc81ee6302c88eb43fd73190f1b036e5f0e6",
		},
		{
			"deposit tx",
			newTestDepositTx(),
			"1d394fb5982739838796ab991da432ddf7123e7bb4b5a3dc92c342458a7d40e1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.tx.Hash()
			require.NoError(t, err)
			assert.Equal(t, tc.hex, hex.EncodeToString(b))
		})
	}
}

func TestTx_Sign(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	addr := crypto.PubkeyToAddress(privKey.PublicKey)

	tx := NewTx()
	require.NoError(t, tx.Sign(0, privKey))

	signers, err := tx.Signers()
	require.NoError(t, err)

	assert.Equal(t, addr, signers[0])
	for i := 1; i < len(signers); i++ {
		assert.Equal(t, NullAddress, signers[i])
	}
}
