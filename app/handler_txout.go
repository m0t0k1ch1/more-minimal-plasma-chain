package app

import "github.com/m0t0k1ch1/more-minimal-plasma-chain/core"

// TODO: can be executed by only the operator
func (p *Plasma) PutTxOutHandler(c *Context) error {
	c.Request().ParseForm()

	txOutPos, err := c.GetTxOutPositionFromPath()
	if err != nil {
		return c.JSONError(err)
	}
	isExited, err := c.GetExitedFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	// BEGIN TXN
	txn := p.db.NewTransaction(true)
	defer txn.Discard()

	if isExited {
		err = p.childChain.ExitTxOut(txn, txOutPos)
	} else {
		err = p.childChain.RestoreTxOut(txn, txOutPos)
	}
	if err != nil {
		if err == core.ErrTxNotFound {
			return c.JSONError(ErrTxNotFound)
		} else if err == core.ErrTxOutNotFound {
			return c.JSONError(ErrTxOutNotFound)
		}
		return c.JSONError(err)
	}

	// COMMIT TXN
	if err := txn.Commit(nil); err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(nil)
}
