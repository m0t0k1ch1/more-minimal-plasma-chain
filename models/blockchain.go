package models

import "sync"

const (
	DefaultBlockNumber = 1
)

type Blockchain struct {
	mu                 *sync.RWMutex
	currentBlockNumber uint64
	chain              map[uint64]*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		mu:                 &sync.RWMutex{},
		currentBlockNumber: DefaultBlockNumber,
		chain:              map[uint64]*Block{},
	}
}

func (bc *Blockchain) GetBlock(num uint64) *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.getBlock(num)
}

func (bc *Blockchain) GetTx(blkNum uint64, txIndex int) *Tx {
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

func (bc *Blockchain) AddBlock(txes []*Tx) (uint64, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	blk, err := NewBlock(txes, bc.currentBlockNumber)
	if err != nil {
		return 0, err
	}

	if _, ok := bc.chain[blk.Number]; ok {
		return 0, ErrBlockAlreadyExists
	}

	bc.chain[blk.Number] = blk
	bc.currentBlockNumber++

	return blk.Number, nil
}

func (bc *Blockchain) getBlock(num uint64) *Block {
	blk, ok := bc.chain[num]
	if !ok {
		return nil
	}

	return blk
}
