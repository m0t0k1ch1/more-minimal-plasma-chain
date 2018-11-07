package core

import (
	"bytes"
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

const (
	DefaultBlockNumber = 1
)

var (
	ErrBlockNotFound = errors.New("block is not found")
	ErrEmptyBlock    = errors.New("block is empty")

	ErrTxNotFound                     = errors.New("tx is not found")
	ErrInvalidTxSignature             = errors.New("tx signature is invalid")
	ErrInvalidTxConfirmationSignature = errors.New("tx confirmation signature is invalid")
	ErrInvalidTxBalance               = errors.New("tx balance is invalid")

	ErrTxInNotFound         = errors.New("txin is not found")
	ErrInvalidTxIn          = errors.New("txin is invalid")
	ErrNullTxInConfirmation = errors.New("null txin cannot be confirmed")

	ErrTxOutAlreadySpent = errors.New("txout is already spent")
)

type ChildChain struct {
	mu           *sync.RWMutex
	currentBlock *types.Block
	chain        map[string]*types.Block // key: blkNum
}

func NewChildChain() (*ChildChain, error) {
	blk, err := types.NewBlock(nil, big.NewInt(DefaultBlockNumber))
	if err != nil {
		return nil, err
	}

	return &ChildChain{
		mu:           &sync.RWMutex{},
		currentBlock: blk,
		chain:        map[string]*types.Block{},
	}, nil
}

func (cc *ChildChain) CurrentBlockNumber() *big.Int {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	return cc.currentBlockNumber()
}

func (cc *ChildChain) GetBlock(blkNum *big.Int) (*types.Block, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	if !cc.isExistBlock(blkNum) {
		return nil, ErrBlockNotFound
	}

	return cc.getBlock(blkNum), nil
}

func (cc *ChildChain) AddBlock(signer *types.Account) (*big.Int, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	blk := cc.currentBlock

	// check block validity
	if len(blk.Txes) == 0 {
		return nil, ErrEmptyBlock
	}

	// sign block
	if err := blk.Sign(signer); err != nil {
		return nil, err
	}

	// add block
	cc.addBlock(blk)

	// reset current block
	blkNext, err := types.NewBlock(nil, cc.newNextBlockNumber())
	if err != nil {
		return nil, err
	}
	cc.currentBlock = blkNext

	return blk.Number, nil
}

func (cc *ChildChain) AddDepositBlock(ownerAddr common.Address, amount *big.Int, signer *types.Account) (*big.Int, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// create deposit tx
	tx := types.NewTx()
	txOut := types.NewTxOut(ownerAddr, amount)
	if err := tx.SetOutput(big.NewInt(0), txOut); err != nil {
		return nil, err
	}

	// create deposit block
	blk, err := types.NewBlock([]*types.Tx{tx}, cc.newCurrentBlockNumber())
	if err != nil {
		return nil, err
	}

	// sign deposit block
	if err := blk.Sign(signer); err != nil {
		return nil, err
	}

	// add deposit block
	cc.addBlock(blk)

	// increment current block number
	cc.incrementBlockNumber()

	return blk.Number, nil
}

func (cc *ChildChain) GetTx(txPos types.Position) (*types.Tx, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	blkNum, txIndex := types.ParseTxPosition(txPos)

	if !cc.isExistTx(blkNum, txIndex) {
		return nil, ErrTxNotFound
	}

	return cc.getTx(blkNum, txIndex), nil
}

func (cc *ChildChain) GetTxProof(txPos types.Position) ([]byte, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	blkNum, txIndex := types.ParseTxPosition(txPos)

	blk := cc.getBlock(blkNum)

	// build tx Merkle tree
	tree, err := blk.MerkleTree()
	if err != nil {
		return nil, err
	}

	// create tx proof
	return tree.CreateMembershipProof(txIndex.Uint64())
}

func (cc *ChildChain) AddTxToMempool(tx *types.Tx) (types.Position, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// validate tx
	if err := cc.validateTx(tx); err != nil {
		return types.Position{}, err
	}

	// add tx to mempool
	if err := cc.addTxToMempool(tx); err != nil {
		return types.Position{}, err
	}

	return types.TxPosition(cc.currentBlock.Number, cc.currentBlock.LastTxIndex()), nil
}

func (cc *ChildChain) ConfirmTx(txInPos types.Position, confSig types.Signature) error {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	blkNum, txIndex, iIndex := types.ParseTxElementPosition(txInPos)

	// check tx existence
	if !cc.isExistTx(blkNum, txIndex) {
		return ErrTxNotFound
	}

	tx := cc.getTx(blkNum, txIndex)

	// check txin existence
	if !tx.IsExistInput(iIndex) {
		return ErrTxInNotFound
	}

	txIn := tx.GetInput(iIndex)

	// check txin validity
	if txIn.IsNull() {
		return ErrNullTxInConfirmation
	}

	inTxOut := cc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

	// verify confirmation signature
	h, err := tx.ConfirmationHash()
	if err != nil {
		return err
	}
	signerAddr, err := confSig.SignerAddress(h)
	if err != nil {
		return ErrInvalidTxConfirmationSignature
	}
	if !bytes.Equal(signerAddr.Bytes(), inTxOut.OwnerAddress.Bytes()) {
		return ErrInvalidTxConfirmationSignature
	}

	// update confirmation signature
	if err := cc.setConfirmationSignature(blkNum, txIndex, iIndex, confSig); err != nil {
		return err
	}

	return nil
}

func (cc *ChildChain) currentBlockNumber() *big.Int {
	return cc.currentBlock.Number
}

func (cc *ChildChain) newCurrentBlockNumber() *big.Int {
	return new(big.Int).Set(cc.currentBlockNumber())
}

func (cc *ChildChain) newNextBlockNumber() *big.Int {
	return new(big.Int).Add(cc.currentBlockNumber(), big.NewInt(1))
}

func (cc *ChildChain) incrementBlockNumber() {
	cc.currentBlockNumber().Add(cc.currentBlockNumber(), big.NewInt(1))
}

func (cc *ChildChain) getBlock(blkNum *big.Int) *types.Block {
	return cc.chain[blkNum.String()]
}

func (cc *ChildChain) isExistBlock(blkNum *big.Int) bool {
	_, ok := cc.chain[blkNum.String()]
	return ok
}

func (cc *ChildChain) addBlock(blk *types.Block) {
	cc.chain[blk.Number.String()] = blk
}

func (cc *ChildChain) addTxToMempool(tx *types.Tx) error {
	for _, txIn := range tx.Inputs {
		if txIn.IsNull() {
			continue
		}

		// spend utxo
		if err := cc.spendUTXO(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex); err != nil {
			return err
		}
	}

	// add tx to current block
	if err := cc.currentBlock.AddTx(tx); err != nil {
		return err
	}

	return nil
}

func (cc *ChildChain) getTx(blkNum, txIndex *big.Int) *types.Tx {
	return cc.chain[blkNum.String()].Txes[txIndex.Int64()]
}

func (cc *ChildChain) isExistTx(blkNum, txIndex *big.Int) bool {
	blk, ok := cc.chain[blkNum.String()]
	if !ok {
		return false
	}

	return blk.IsExistTx(txIndex)
}

func (cc *ChildChain) validateTx(tx *types.Tx) error {
	nullTxInNum := 0
	iAmount, oAmount := big.NewInt(0), big.NewInt(0)

	for _, txOut := range tx.Outputs {
		oAmount.Add(oAmount, txOut.Amount)
	}

	for i, txIn := range tx.Inputs {
		// check spending txout existence
		if !cc.isExistTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex) {
			if txIn.IsNull() {
				nullTxInNum++
				continue
			}
			return ErrInvalidTxIn
		}

		inTxOut := cc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

		// check double spent
		if inTxOut.IsSpent {
			return ErrTxOutAlreadySpent
		}

		// verify signature
		signerAddr, err := tx.SignerAddress(big.NewInt(int64(i)))
		if err != nil {
			return ErrInvalidTxSignature
		}
		if txIn.Signature == types.NullSignature ||
			!bytes.Equal(signerAddr.Bytes(), inTxOut.OwnerAddress.Bytes()) {
			return ErrInvalidTxSignature
		}

		iAmount.Add(iAmount, inTxOut.Amount)
	}

	// check txins validity
	if nullTxInNum == len(tx.Inputs) {
		return ErrInvalidTxIn
	}

	// check in/out balance
	if iAmount.Cmp(oAmount) < 0 {
		return ErrInvalidTxBalance
	}

	return nil
}

func (cc *ChildChain) getTxOut(blkNum, txIndex, oIndex *big.Int) *types.TxOut {
	return cc.getTx(blkNum, txIndex).GetOutput(oIndex)
}

func (cc *ChildChain) isExistTxOut(blkNum, txIndex, oIndex *big.Int) bool {
	if !cc.isExistTx(blkNum, txIndex) {
		return false
	}

	return cc.getTx(blkNum, txIndex).IsExistOutput(oIndex)
}

func (cc *ChildChain) spendUTXO(blkNum, txIndex, oIndex *big.Int) error {
	return cc.getTx(blkNum, txIndex).SpendOutput(oIndex)
}

func (cc *ChildChain) setConfirmationSignature(blkNum, txIndex, iIndex *big.Int, confSig types.Signature) error {
	return cc.getTx(blkNum, txIndex).SetConfirmationSignature(iIndex, confSig)
}
