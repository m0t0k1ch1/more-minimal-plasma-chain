package models

import (
	"bytes"
	"crypto/ecdsa"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type Block struct {
	Txes      []*Tx
	Number    uint
	Signature []byte
}

func NewBlock(num uint) *Block {
	return &Block{
		Txes:      nil,
		Number:    num,
		Signature: NullSignature,
	}
}

// implements RLP Encoder interface
// ref. https://godoc.org/github.com/ethereum/go-ethereum/rlp#Encoder
func (blk *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		blk.Txes, blk.Number,
	})
}

func (blk *Block) Hash() ([]byte, error) {
	b, err := rlp.EncodeToBytes(blk)
	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(b), nil
}

func (blk *Block) Sign(privKey *ecdsa.PrivateKey) error {
	hashBytes, err := blk.Hash()
	if err != nil {
		return err
	}

	sig, err := crypto.Sign(hashBytes, privKey)
	if err != nil {
		return err
	}
	blk.Signature = sig

	return nil
}

func (blk *Block) Signer() (common.Address, error) {
	hashBytes, err := blk.Hash()
	if err != nil {
		return common.Address{}, err
	}

	if bytes.Equal(blk.Signature, NullSignature) {
		return NullAddress, nil
	}

	pubKey, err := crypto.SigToPub(hashBytes, blk.Signature)
	if err != nil {
		return common.Address{}, err
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}
