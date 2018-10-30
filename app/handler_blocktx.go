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

	tx, err := cc.blockchain.GetTx(blkNum, txIndex)
	if err != nil {
		if err == core.ErrTxNotFound {
			return c.JSONError(ErrTxNotFound)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(tx)
}

func (cc *ChildChain) PutBlockTxHandler(c *Context) error {
	c.Request().ParseForm()

	blkNum, err := c.GetBlockNumberFromPath()
	if err != nil {
		return c.JSONError(err)
	}
	txIndex, err := c.GetTxIndexFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	iIndex, err := c.GetIndexFromForm()
	if err != nil {
		return c.JSONError(err)
	}
	confSig, err := c.GetConfirmationSignatureFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	if err := cc.blockchain.SetConfirmationSignature(blkNum, txIndex, iIndex, confSig); err != nil {
		if err == core.ErrTxInNotFound {
			return c.JSONError(ErrTxInNotFound)
		} else if err == core.ErrInvalidTxConfirmationSignature {
			return c.JSONError(ErrInvalidTxConfirmationSignature)
		}
		return c.JSONError(err)
	}

	return nil
}
