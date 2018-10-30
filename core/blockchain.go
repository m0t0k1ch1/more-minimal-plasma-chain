package core

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

const (
	DefaultBlockNumber = 1
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

	// TODO: validate tx

	bc.currentBlock.Txes = append(bc.currentBlock.Txes, tx)

	return nil
}

func (bc *Blockchain) GetTx(blkNum uint64, txIndex int) *types.Tx {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	blk := bc.getBlock(blkNum)
	if blk == nil {
		return nil
	}

	if txIndex >= len(blk.Txes) {
		return nil
	}

	return blk.Txes[txIndex]
}

func (bc *Blockchain) getBlock(blkNum uint64) *types.Block {
	blk, ok := bc.chain[blkNum]
	if !ok {
		return nil
	}

	return blk
}
