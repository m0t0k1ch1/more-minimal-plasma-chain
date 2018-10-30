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
	ErrInvalidTxInput      = errors.New("tx input is invalid")
	ErrTxInputAlreadySpent = errors.New("tx input is already spent")
	ErrInvalidTxSignature  = errors.New("tx signature is invalid")
	ErrInvalidTxBalance    = errors.New("tx balance is invalid")
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

func (bc *Blockchain) GetBlock(blkNum uint64) *types.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.getBlock(blkNum)
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
		bc.chain[txIn.BlockNumber].Txes[txIn.TxIndex].Spents[txIn.OutputIndex] = true
	}
	bc.currentBlock.Txes = append(bc.currentBlock.Txes, tx)

	return nil
}

func (bc *Blockchain) GetTx(blkNum uint64, txIndex uint64) *types.Tx {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.getTx(blkNum, txIndex)
}

func (bc *Blockchain) validateTx(tx *types.Tx) error {
	signerAddrs, err := tx.SignerAddresses()
	if err != nil {
		return err
	}

	iAmount, oAmount := uint64(0), uint64(0)

	for _, txOut := range tx.Outputs {
		oAmount += txOut.Amount
	}

	for i, txIn := range tx.Inputs {
		if txIn.BlockNumber == 0 {
			continue
		}

		inTx := bc.getTx(txIn.BlockNumber, txIn.TxIndex)
		if inTx == nil ||
			txIn.OutputIndex >= uint64(len(inTx.Outputs)) ||
			txIn.OutputIndex >= uint64(len(inTx.Spents)) {
			return ErrInvalidTxInput
		}
		if inTx.Spents[txIn.OutputIndex] {
			return ErrTxInputAlreadySpent
		}

		inTxOut := inTx.Outputs[txIn.OutputIndex]
		iAmount += inTxOut.Amount

		if tx.Signatures[i] == types.NullSignature ||
			!bytes.Equal(inTxOut.OwnerAddress.Bytes(), signerAddrs[i].Bytes()) {
			return ErrInvalidTxSignature
		}
	}

	if !tx.IsDeposit() && iAmount < oAmount {
		return ErrInvalidTxBalance
	}

	return nil
}

func (bc *Blockchain) getBlock(blkNum uint64) *types.Block {
	blk, ok := bc.chain[blkNum]
	if !ok {
		return nil
	}

	return blk
}

func (bc *Blockchain) getTx(blkNum uint64, txIndex uint64) *types.Tx {
	blk := bc.getBlock(blkNum)
	if blk == nil {
		return nil
	}

	if txIndex >= uint64(len(blk.Txes)) {
		return nil
	}

	return blk.Txes[txIndex]
}
