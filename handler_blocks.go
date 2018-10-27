package main

import (
	"strconv"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/models"
)

func (cc *ChildChain) GetBlockHandler(c *Context) error {
	numStr := c.Param("num")

	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		return err
	}

	return c.JSONSuccess(cc.getBlock(num))
}

func (cc *ChildChain) PostBlockHandler(c *Context) error {
	cc.addBlock()
	return c.JSONSuccess(nil)
}

func (cc *ChildChain) getBlock(num uint64) *models.Block {
	return cc.blocks[num]
}

func (cc *ChildChain) addBlock() {
	blk := models.NewBlock(cc.currentBlockNumber)
	cc.blocks[blk.Number] = blk
	cc.currentBlockNumber++
}
