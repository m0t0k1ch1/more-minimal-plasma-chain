package main

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/models"
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

func (cc *ChildChain) PostDepositBlockHandler(c *Context) error {
	ownerStr := c.FormValue("owner")
	amountStr := c.FormValue("amount")

	if !common.IsHexAddress(ownerStr) {
		return c.JSONError(ErrInvalidAddressHex)
	}
	owner := common.HexToAddress(ownerStr)

	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return c.JSONError(err)
	}

	tx := models.NewTx()
	tx.Outputs[0] = models.NewTxOut(owner, amount)

	num, err := cc.blockchain.AddBlock([]*models.Tx{tx})
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
