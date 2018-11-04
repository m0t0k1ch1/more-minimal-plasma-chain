package types

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	merkle "github.com/m0t0k1ch1/fixed-merkle"
)

const (
	MaxBlockTxesNum = 99999
)

var (
	ErrInvalidTxIndex           = errors.New("tx index is invalid")
	ErrBlockTxesNumExceedsLimit = errors.New("block txes num exceeds the limit")
)

type LightBlock struct {
	TxHashes  []common.Hash
	Number    *big.Int
	Signature Signature
}

func (lblk *LightBlock) GetTxHash(txIndex *big.Int) common.Hash {
	if !lblk.IsExistTxHash(txIndex) {
		return NullHash
	}

	return lblk.TxHashes[txIndex.Uint64()]
}

func (lblk *LightBlock) IsExistTxHash(txIndex *big.Int) bool {
	return txIndex.Cmp(big.NewInt(int64(len(lblk.TxHashes)))) < 0
}

type Block struct {
	Txes      []*Tx     `json:"txes"`
	Number    *big.Int  `json:"blknum"`
	Signature Signature `json:"sig"`
}

func NewBlock(txes []*Tx, blkNum *big.Int) (*Block, error) {
	if len(txes) >= MaxBlockTxesNum {
		return nil, ErrBlockTxesNumExceedsLimit
	}

	return &Block{
		Txes:      txes,
		Number:    blkNum,
		Signature: NullSignature,
	}, nil
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

func (blk *Block) Hash() (common.Hash, error) {
	b, err := blk.Encode()
	if err != nil {
		return NullHash, err
	}

	return common.BytesToHash(crypto.Keccak256(b)), nil
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

func (blk *Block) Root() (common.Hash, error) {
	rootHash := common.Hash{}

	tree, err := blk.MerkleTree()
	if err != nil {
		return rootHash, err
	}

	copy(rootHash[:], tree.Root().Bytes()[:])

	return rootHash, nil
}

func (blk *Block) AddTx(tx *Tx) error {
	if len(blk.Txes) >= MaxBlockTxesNum {
		return ErrBlockTxesNumExceedsLimit
	}

	blk.Txes = append(blk.Txes, tx)

	return nil
}

func (blk *Block) Sign(signer *Account) error {
	h, err := blk.Hash()
	if err != nil {
		return err
	}

	sigBytes, err := signer.Sign(h)
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
	h, err := blk.Hash()
	if err != nil {
		return NullAddress, err
	}

	if bytes.Equal(blk.Signature.Bytes(), NullSignature.Bytes()) {
		return NullAddress, nil
	}

	return blk.Signature.SignerAddress(h)
}

func (blk *Block) Lighten() (*LightBlock, error) {
	lblk := &LightBlock{
		TxHashes:  make([]common.Hash, len(blk.Txes)),
		Number:    blk.Number,
		Signature: blk.Signature,
	}

	for i, tx := range blk.Txes {
		txHash, err := tx.Hash()
		if err != nil {
			return nil, err
		}

		lblk.TxHashes[i] = txHash
	}

	return lblk, nil
}
