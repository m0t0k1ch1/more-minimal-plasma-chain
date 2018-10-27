package main

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
			NullTx,
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
			assert.Equal(t, hex.EncodeToString(b), tc.hex)
		})
	}
}
