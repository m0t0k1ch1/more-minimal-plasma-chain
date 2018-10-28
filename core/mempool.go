package core

import (
	"errors"
	"sync"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

var (
	ErrMempoolFull = errors.New("mempool is full")
)

type Mempool struct {
	mu     *sync.RWMutex
	offset int
	pool   []*types.Tx
}

func NewMempool(size int) *Mempool {
	mp := &Mempool{
		mu:     &sync.RWMutex{},
		offset: 0,
		pool:   make([]*types.Tx, size),
	}

	return mp
}

func (mp *Mempool) Add(tx *types.Tx) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.offset >= len(mp.pool) {
		return ErrMempoolFull
	}

	mp.pool[mp.offset] = tx
	mp.offset++

	return nil
}

func (mp *Mempool) Extract() []*types.Tx {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	txes := make([]*types.Tx, mp.offset)
	for i := 0; i < mp.offset; i++ {
		txes[i] = mp.pool[i]
		mp.pool[i] = nil
	}
	mp.offset = 0

	return txes
}
