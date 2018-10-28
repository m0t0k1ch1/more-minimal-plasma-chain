package core

import (
	"sync"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

const (
	DefaultBlockNumber = 1
)

type Blockchain struct {
	mu                 *sync.RWMutex
	currentBlockNumber uint64
	chain              map[uint64]*types.Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		mu:                 &sync.RWMutex{},
		currentBlockNumber: DefaultBlockNumber,
		chain:              map[uint64]*types.Block{},
	}
}

func (bc *Blockchain) GetBlock(blkNum uint64) *types.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.getBlock(blkNum)
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

func (bc *Blockchain) AddBlock(txes []*types.Tx) (uint64, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	blk, err := types.NewBlock(txes, bc.currentBlockNumber)
	if err != nil {
		return 0, err
	}

	bc.chain[blk.Number] = blk
	bc.currentBlockNumber++

	return blk.Number, nil
}

func (bc *Blockchain) getBlock(blkNum uint64) *types.Block {
	blk, ok := bc.chain[blkNum]
	if !ok {
		return nil
	}

	return blk
}
