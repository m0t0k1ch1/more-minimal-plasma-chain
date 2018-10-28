package app

func (cc *ChildChain) GetBlockTxHandler(c *Context) error {
	blkNum, err := c.GetBlockNumberFromPath()
	if err != nil {
		return c.JSONError(err)
	}
	txIndex, err := c.GetTxIndexFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	tx := cc.blockchain.GetTx(blkNum, txIndex)
	if tx == nil {
		return c.JSONError(ErrTxNotFound)
	}

	return c.JSONSuccess(tx)
}

func (cc *ChildChain) PutBlockTxHandler(c *Context) error {
	return nil
}
