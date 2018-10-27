package main

import (
	"strconv"
)

func (cc *ChildChain) GetBlockHandler(c *Context) error {
	numStr := c.Param("num")

	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		return c.JSONError(err)
	}

	return cc.getBlockHandler(c, num)
}

func (cc *ChildChain) PostBlockHandler(c *Context) error {
	txes := cc.mempool.Extract()

	num, err := cc.blockchain.AddBlock(txes)
	if err != nil {
		return c.JSONError(err)
	}

	return cc.getBlockHandler(c, num)
}

func (cc *ChildChain) getBlockHandler(c *Context, num uint64) error {
	blk := cc.blockchain.GetBlock(num)
	if blk == nil {
		return c.JSONError(ErrBlockNotFound)
	}

	blkSummary, err := blk.Summary()
	if err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(blkSummary)
}
