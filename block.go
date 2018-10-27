package main

import (
	"io"

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
