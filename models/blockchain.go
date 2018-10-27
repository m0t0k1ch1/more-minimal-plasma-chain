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

	blk, ok := bc.chain[num]
	if !ok {
		return nil
	}

	return blk
}

func (bc *Blockchain) AddBlock(txes []*Tx) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	blk, err := NewBlock(txes, bc.currentBlockNumber)
	if err != nil {
		return err
	}

	if _, ok := bc.chain[blk.Number]; ok {
		return ErrBlockAlreadyExists
	}

	bc.chain[blk.Number] = blk
	bc.currentBlockNumber++

	return nil
}
