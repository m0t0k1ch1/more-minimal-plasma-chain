package types

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	merkle "github.com/m0t0k1ch1/fixed-merkle"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

const (
	MaxBlockTxesNum = 99999
)

var (
	ErrInvalidTxIndex           = errors.New("tx index is invalid")
	ErrBlockTxesNumExceedsLimit = errors.New("block txes num exceeds the limit")
)

type BlockHeader struct {
	Number    uint64
	Signature Signature
}

type Block struct {
	*BlockHeader
	Txes []*Tx `json:"txes"`
}

func NewBlock(txes []*Tx, blkNum uint64) (*Block, error) {
	if len(txes) >= MaxBlockTxesNum {
		return nil, ErrBlockTxesNumExceedsLimit
	}

	return &Block{
		BlockHeader: &BlockHeader{
			Number:    blkNum,
			Signature: NullSignature,
		},
		Txes: txes,
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
	blkRootHash := common.Hash{}

	tree, err := blk.MerkleTree()
	if err != nil {
		return blkRootHash, err
	}

	return utils.BytesToHash(tree.Root().Bytes()), nil
}

func (blk *Block) IsExistTx(txIndex uint64) bool {
	return txIndex < uint64(len(blk.Txes))
}

func (blk *Block) GetTx(txIndex uint64) *Tx {
	if !blk.IsExistTx(txIndex) {
		return nil
	}

	return blk.Txes[txIndex]
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
	sig, err := BytesToSignature(sigBytes)
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
