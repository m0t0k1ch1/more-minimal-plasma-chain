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

func (cc *ChildChain) PostBlockHandler(c *Context) error {
	txes := cc.mempool.Extract()

	if err := cc.blockchain.AddBlock(txes); err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(nil)
}
