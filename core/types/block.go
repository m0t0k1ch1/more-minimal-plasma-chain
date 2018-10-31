package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	merkle "github.com/m0t0k1ch1/fixed-merkle"
)

type Block struct {
	Txes      []*Tx     `json:"txes"`
	Number    uint64    `json:"blknum"`
	Signature Signature `json:"sig"`
}

func NewBlock(txes []*Tx, blkNum uint64) *Block {
	return &Block{
		Txes:      txes,
		Number:    blkNum,
		Signature: NullSignature,
	}
}

func (blk *Block) Encode() ([]byte, error) {
	txCores := make([]interface{}, len(blk.Txes))
	for i, tx := range blk.Txes {
		txCores[i] = []interface{}{
			tx.inputCores(), tx.outputCores(), tx.signatures(),
		}
	}

	return rlp.EncodeToBytes([]interface{}{
		txCores, blk.Number,
	})
}

func (blk *Block) Hash() ([]byte, error) {
	b, err := blk.Encode()
	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(b), nil
}

func (blk *Block) MerkleTree() (*merkle.Tree, error) {
	leaves := make([][]byte, len(blk.Txes))
	for i, tx := range blk.Txes {
		leaf, err := tx.MerkleLeaf()
		if err != nil {
			return nil, err
		}

		leaves[i] = leaf
	}

	return merkle.NewTree(merkleConfig(), leaves)
}

func (blk *Block) Root() ([]byte, error) {
	tree, err := blk.MerkleTree()
	if err != nil {
		return nil, err
	}

	return tree.Root().Bytes(), nil
}

func (blk *Block) Sign(signer *Account) error {
	hashBytes, err := blk.Hash()
	if err != nil {
		return err
	}

	sigBytes, err := signer.Sign(hashBytes)
	if err != nil {
		return err
	}

	sig, err := NewSignatureFromBytes(sigBytes)
	if err != nil {
		return err
	}

	blk.Signature = sig

	return nil
}

func (blk *Block) SignerAddress() (common.Address, error) {
	hashBytes, err := blk.Hash()
	if err != nil {
		return NullAddress, err
	}

	if bytes.Equal(blk.Signature.Bytes(), NullSignature.Bytes()) {
		return NullAddress, nil
	}

	return blk.Signature.SignerAddress(hashBytes)
}
