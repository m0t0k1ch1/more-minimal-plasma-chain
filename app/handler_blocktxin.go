package app

import "github.com/m0t0k1ch1/more-minimal-plasma-chain/core"

func (cc *ChildChain) PutBlockTxInHandler(c *Context) error {
	c.Request().ParseForm()

	blkNum, err := c.GetBlockNumberFromPath()
	if err != nil {
		return c.JSONError(err)
	}
	txIndex, err := c.GetTxIndexFromPath()
	if err != nil {
		return c.JSONError(err)
	}
	iIndex, err := c.GetInputIndexFromPath()
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
		} else if err == core.ErrDepositTxInConfirmation {
			return c.JSONError(ErrDepositTxInConfirmation)
		} else if err == core.ErrInvalidTxConfirmationSignature {
			return c.JSONError(ErrInvalidTxConfirmationSignature)
		}
		return c.JSONError(err)
	}

	return cc.getBlockTxInHandler(c, blkNum, txIndex, iIndex)
}

func (cc *ChildChain) getBlockTxInHandler(c *Context, blkNum, txIndex, iIndex uint64) error {
	txIn, err := cc.blockchain.GetTxIn(blkNum, txIndex, iIndex)
	if err != nil {
		if err == core.ErrTxInNotFound {
			return c.JSONError(ErrTxInNotFound)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(txIn)
}
