package app

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
)

func (p *Plasma) PutTxInHandler(c *Context) error {
	c.Request().ParseForm()

	txInPos, err := c.GetTxInPositionFromPath()
	if err != nil {
		return c.JSONError(err)
	}
	confSig, err := c.GetConfirmationSignatureFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	// BEGIN TXN
	txn := p.db.NewTransaction(true)
	defer txn.Discard()

	if err := p.childChain.ConfirmTx(txn, txInPos, confSig); err != nil {
		if err == core.ErrTxNotFound {
			return c.JSONError(ErrTxNotFound)
		} else if err == core.ErrTxInNotFound {
			return c.JSONError(ErrTxInNotFound)
		} else if err == core.ErrNullTxInConfirmation {
			return c.JSONError(ErrNullTxInConfirmation)
		} else if err == core.ErrInvalidTxIn {
			return c.JSONError(ErrInvalidTxIn)
		} else if err == core.ErrInvalidTxConfirmationSignature {
			return c.JSONError(ErrInvalidTxConfirmationSignature)
		}
		return c.JSONError(err)
	}

	// COMMIT TXN
	if err := txn.Commit(nil); err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(nil)
}
