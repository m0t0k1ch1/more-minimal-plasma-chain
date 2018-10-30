package core

import (
	"bytes"
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

const (
	DefaultBlockNumber = 1
)

var (
	ErrBlockNotFound = errors.New("block is not found")

	ErrTxNotFound                     = errors.New("tx is not found")
	ErrInvalidTxSignature             = errors.New("tx signature is invalid")
	ErrInvalidTxConfirmationSignature = errors.New("tx confirmation signature is invalid")
	ErrInvalidTxBalance               = errors.New("tx balance is invalid")

	ErrTxInNotFound = errors.New("txin is not found")
	ErrInvalidTxIn  = errors.New("txin is invalid")

	ErrTxOutAlreadySpent = errors.New("txout is already spent")
)

type Blockchain struct {
	mu           *sync.RWMutex
	currentBlock *types.Block
	chain        map[uint64]*types.Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		mu:           &sync.RWMutex{},
		currentBlock: types.NewBlock(nil, DefaultBlockNumber),
		chain:        map[uint64]*types.Block{},
	}
}

func (bc *Blockchain) CurrentBlockNumber() uint64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.currentBlock.Number
}

func (bc *Blockchain) GetBlock(blkNum uint64) (*types.Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if !bc.isExistBlock(blkNum) {
		return nil, ErrBlockNotFound
	}

	return bc.getBlock(blkNum), nil
}

func (bc *Blockchain) AddBlock(signer *types.Account) (uint64, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if err := bc.currentBlock.Sign(signer); err != nil {
		return 0, err
	}

	blkNum := bc.currentBlock.Number
	bc.chain[blkNum] = bc.currentBlock
	bc.currentBlock = types.NewBlock(nil, blkNum+1)

	return blkNum, nil
}

func (bc *Blockchain) AddDepositBlock(ownerAddr common.Address, amount uint64, signer *types.Account) (uint64, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	tx := types.NewTx()
	tx.Outputs[0] = types.NewTxOut(ownerAddr, amount)

	blk := types.NewBlock([]*types.Tx{tx}, bc.currentBlock.Number)
	if err := blk.Sign(signer); err != nil {
		return 0, err
	}

	bc.chain[blk.Number] = blk
	bc.currentBlock.Number++

	return blk.Number, nil
}

func (bc *Blockchain) GetTx(blkNum, txIndex uint64) (*types.Tx, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if !bc.isExistTx(blkNum, txIndex) {
		return nil, ErrTxNotFound
	}

	return bc.getTx(blkNum, txIndex), nil
}

func (bc *Blockchain) AddTx(tx *types.Tx) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if err := bc.validateTx(tx); err != nil {
		return err
	}

	for _, txIn := range tx.Inputs {
		if txIn.BlockNumber == 0 {
			continue
		}
		bc.chain[txIn.BlockNumber].Txes[txIn.TxIndex].Outputs[txIn.OutputIndex].Spent()
	}
	bc.currentBlock.Txes = append(bc.currentBlock.Txes, tx)

	return nil
}

func (bc *Blockchain) GetTxIn(blkNum, txIndex, iIndex uint64) (*types.TxIn, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if !bc.isExistTxIn(blkNum, txIndex, iIndex) {
		return nil, ErrTxInNotFound
	}

	return bc.getTxIn(blkNum, txIndex, iIndex), nil
}

func (bc *Blockchain) SetConfirmationSignature(blkNum, txIndex, iIndex uint64, confSig types.Signature) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if !bc.isExistTxIn(blkNum, txIndex, iIndex) {
		return ErrTxInNotFound
	}

	tx := bc.getTx(blkNum, txIndex)
	txIn := tx.Inputs[iIndex]
	inTxOut := bc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

	confHashBytes, err := tx.ConfirmationHash()
	if err != nil {
		return err
	}
	confSigner, err := confSig.SignerAddress(confHashBytes)
	if err != nil {
		return ErrInvalidTxConfirmationSignature
	}
	if confSigner.Hex() != inTxOut.OwnerAddress.Hex() {
		return ErrInvalidTxConfirmationSignature
	}

	bc.chain[blkNum].Txes[txIndex].Inputs[iIndex].ConfirmationSignature = confSig

	return nil
}

func (bc *Blockchain) validateTx(tx *types.Tx) error {
	iAmount, oAmount := uint64(0), uint64(0)

	for _, txOut := range tx.Outputs {
		oAmount += txOut.Amount
	}

	for i, txIn := range tx.Inputs {
		if txIn.BlockNumber == 0 {
			continue
		}

		if !bc.isExistTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex) {
			return ErrInvalidTxIn
		}

		inTxOut := bc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

		if inTxOut.IsSpent {
			return ErrTxOutAlreadySpent
		}

		signerAddr, err := tx.SignerAddress(uint64(i))
		if err != nil {
			return ErrInvalidTxSignature
		}
		if txIn.Signature == types.NullSignature ||
			!bytes.Equal(signerAddr.Bytes(), inTxOut.OwnerAddress.Bytes()) {
			return ErrInvalidTxSignature
		}

		iAmount += inTxOut.Amount
	}

	if !tx.IsDeposit() && iAmount < oAmount {
		return ErrInvalidTxBalance
	}

	return nil
}

func (bc *Blockchain) isExistBlock(blkNum uint64) bool {
	_, ok := bc.chain[blkNum]
	return ok
}

func (bc *Blockchain) getBlock(blkNum uint64) *types.Block {
	return bc.chain[blkNum]
}

func (bc *Blockchain) isExistTx(blkNum, txIndex uint64) bool {
	if !bc.isExistBlock(blkNum) {
		return false
	}

	return txIndex < uint64(len(bc.chain[blkNum].Txes))
}

func (bc *Blockchain) getTx(blkNum, txIndex uint64) *types.Tx {
	return bc.chain[blkNum].Txes[txIndex]
}

func (bc *Blockchain) isExistTxIn(blkNum, txIndex, iIndex uint64) bool {
	if !bc.isExistTx(blkNum, txIndex) {
		return false
	}

	return iIndex < uint64(len(bc.chain[blkNum].Txes[txIndex].Inputs))
}

func (bc *Blockchain) getTxIn(blkNum, txIndex, iIndex uint64) *types.TxIn {
	return bc.chain[blkNum].Txes[txIndex].Inputs[iIndex]
}

func (bc *Blockchain) isExistTxOut(blkNum, txIndex, oIndex uint64) bool {
	if !bc.isExistTx(blkNum, txIndex) {
		return false
	}

	return oIndex < uint64(len(bc.chain[blkNum].Txes[txIndex].Outputs))
}

func (bc *Blockchain) getTxOut(blkNum, txIndex, oIndex uint64) *types.TxOut {
	return bc.chain[blkNum].Txes[txIndex].Outputs[oIndex]
}
