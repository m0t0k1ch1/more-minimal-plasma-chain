package models

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestNullTx(t *testing.T) *Tx {
	return NewTx()
}

func newTestDepositTx(t *testing.T) *Tx {
	tx := NewTx()
	tx.Outputs[0] = NewTxOut(
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		1,
	)
	return tx
}

func TestTx_Hash(t *testing.T) {
	type output struct {
		hashHex string
		err     error
	}
	testCases := []struct {
		name string
		tx   *Tx
		out  output
	}{
		{
			"null tx",
			newTestNullTx(t),
			output{
				"0xc758c57a2f76021ff85aa579dc03dc81ee6302c88eb43fd73190f1b036e5f0e6",
				nil,
			},
		},
		{
			"deposit tx",
			newTestDepositTx(t),
			output{
				"0xcfbf8d16cf5cd7a8f3a812ad415a163b0d112c08b061ce36f6291dae81a97f8e",
				nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx, out := tc.tx, tc.out

			hashBytes, err := tx.Hash()
			if out.err != nil {
				assert.EqualError(t, err, out.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, out.hashHex, hexutil.Encode(hashBytes))
			}
		})
	}
}

func TestTx_ConfirmationHash(t *testing.T) {
	type output struct {
		confHashHex string
		err         error
	}
	testCases := []struct {
		name string
		tx   *Tx
		out  output
	}{
		{
			"null tx",
			newTestNullTx(t),
			output{
				"0xe8bd8330f57eeb0fe2d8afd9f5bcc53ebcddd1378054a6bbdb224e48ddd50ca1",
				nil,
			},
		},
		{
			"deposit tx",
			newTestDepositTx(t),
			output{
				"0x26a410eed92c10d92d54d131f9d77896d7b0fd4d4582b8650e6db5ecf112e05f",
				nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx, out := tc.tx, tc.out

			confHashBytes, err := tx.ConfirmationHash()
			if out.err != nil {
				assert.EqualError(t, err, out.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, out.confHashHex, hexutil.Encode(confHashBytes))
			}
		})
	}
}

func TestTx_Sign(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	tx := newTestNullTx(t)

	// sign
	require.NoError(t, tx.Sign(0, privKey))

	// verify
	signers, err := tx.Signers()
	require.NoError(t, err)
	assert.Equal(t, crypto.PubkeyToAddress(privKey.PublicKey), signers[0])
	for i := 1; i < len(signers); i++ {
		assert.Equal(t, nullAddress, signers[i])
	}
}
