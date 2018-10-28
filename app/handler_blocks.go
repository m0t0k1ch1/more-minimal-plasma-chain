package app

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

func (cc *ChildChain) GetBlockHandler(c *Context) error {
	blkNum, err := c.ParamBlockNumber()
	if err != nil {
		return c.JSONError(err)
	}

	return cc.getBlockHandler(c, blkNum)
}

func (cc *ChildChain) PostBlockHandler(c *Context) error {
	txes := cc.mempool.Extract()

	blkNum, err := cc.blockchain.AddBlock(txes)
	if err != nil {
		return c.JSONError(err)
	}

	return cc.getBlockHandler(c, blkNum)
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

	tx := types.NewTx()
	tx.Outputs[0] = types.NewTxOut(owner, amount)

	blkNum, err := cc.blockchain.AddBlock([]*types.Tx{tx})
	if err != nil {
		return c.JSONError(err)
	}

	return cc.getBlockHandler(c, blkNum)
}

func (cc *ChildChain) GetBlockTxHandler(c *Context) error {
	blkNum, err := c.ParamBlockNumber()
	if err != nil {
		return c.JSONError(err)
	}
	txIndex, err := c.ParamTxIndex()
	if err != nil {
		return c.JSONError(err)
	}

	tx := cc.blockchain.GetTx(blkNum, txIndex)
	if tx == nil {
		return c.JSONError(ErrTxNotFound)
	}

	return c.JSONSuccess(tx)
}

func (cc *ChildChain) getBlockHandler(c *Context, blkNum uint64) error {
	blk := cc.blockchain.GetBlock(blkNum)
	if blk == nil {
		return c.JSONError(ErrBlockNotFound)
	}

	blkSummary, err := blk.Summary()
	if err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(blkSummary)
}
