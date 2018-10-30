package app

import "github.com/m0t0k1ch1/more-minimal-plasma-chain/core"

func (cc *ChildChain) PostTxHandler(c *Context) error {
	c.Request().ParseForm()

	tx, err := c.GetTxFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	if err := cc.blockchain.AddTx(tx); err != nil {
		if err == core.ErrInvalidTxSignature {
			return c.JSONError(ErrInvalidTxSignature)
		} else if err == core.ErrInvalidTxBalance {
			return c.JSONError(ErrInvalidTxBalance)
		} else if err == core.ErrInvalidTxIn {
			return c.JSONError(ErrInvalidTxIn)
		} else if err == core.ErrTxOutAlreadySpent {
			return c.JSONError(ErrTxOutAlreadySpent)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(tx)
}
