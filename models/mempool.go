package models

import (
	"errors"
	"sync"
)

var (
	ErrMempoolFull = errors.New("mempool is full")
)

type Mempool struct {
	mu     *sync.RWMutex
	offset int
	pool   []*Tx
}

func NewMempool(size int) *Mempool {
	mp := &Mempool{
		mu:     &sync.RWMutex{},
		offset: 0,
		pool:   make([]*Tx, size),
	}

	return mp
}

func (mp *Mempool) Add(tx *Tx) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.offset >= len(mp.pool) {
		return ErrMempoolFull
	}

	mp.pool[mp.offset] = tx
	mp.offset++

	return nil
}

func (mp *Mempool) Extract() []*Tx {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	txes := make([]*Tx, mp.offset)
	for i := 0; i < mp.offset; i++ {
		txes[i] = mp.pool[i]
		mp.pool[i] = nil
	}
	mp.offset = 0

	return txes
}
