package main

import (
	"strconv"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/models"
)

func (cc *ChildChain) GetBlockHandler(c *Context) error {
	numStr := c.Param("num")

	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		return c.JSONError(err)
	}

	blk, err := cc.getBlock(num)
	if err != nil {
		return c.JSONError(err)
	}

	blkSummary, err := blk.Summary()
	if err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(blkSummary)
}

func (cc *ChildChain) PostBlockHandler(c *Context) error {
	if err := cc.addBlock(); err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(nil)
}

func (cc *ChildChain) getBlock(num uint64) (*models.Block, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	blk, ok := cc.blocks[num]
	if !ok {
		return nil, ErrBlockNotFound
	}

	return blk, nil
}

func (cc *ChildChain) addBlock() error {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	txes := cc.mempool.Extract()

	blk, err := models.NewBlock(txes, cc.currentBlockNumber)
	if err != nil {
		return err
	}

	cc.blocks[blk.Number] = blk
	cc.currentBlockNumber++

	return nil
}
