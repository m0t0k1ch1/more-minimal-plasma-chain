package app

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
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

	if err := p.childChain.ConfirmTx(txInPos, confSig); err != nil {
		if err == core.ErrInvalidTxConfirmationSignature {
			return c.JSONError(ErrInvalidTxConfirmationSignature)
		} else if err == core.ErrTxInNotFound {
			return c.JSONError(ErrTxInNotFound)
		} else if err == core.ErrNullTxInConfirmation {
			return c.JSONError(ErrNullTxInConfirmation)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]*types.Position{
		"pos": txInPos,
	})
}
