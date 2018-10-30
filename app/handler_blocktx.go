package app

import "github.com/m0t0k1ch1/more-minimal-plasma-chain/core"

func (cc *ChildChain) GetBlockTxHandler(c *Context) error {
	blkNum, err := c.GetBlockNumberFromPath()
	if err != nil {
		return c.JSONError(err)
	}
	txIndex, err := c.GetTxIndexFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	return cc.getBlockTxHandler(c, blkNum, txIndex)
}

func (cc *ChildChain) getBlockTxHandler(c *Context, blkNum, txIndex uint64) error {
	tx, err := cc.blockchain.GetTx(blkNum, txIndex)
	if err != nil {
		if err == core.ErrTxNotFound {
			return c.JSONError(ErrTxNotFound)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(tx)
}
