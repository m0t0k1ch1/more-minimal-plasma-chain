package app

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

func (p *Plasma) PostTxHandler(c *Context) error {
	c.Request().ParseForm()

	tx, err := c.GetTxFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	// BEGIN TXN
	txn := p.db.NewTransaction(true)
	defer txn.Discard()

	if err := p.childChain.AddTxToMempool(txn, tx); err != nil {
		if err == core.ErrMempoolFull {
			return c.JSONError(ErrMempoolFull)
		} else if err == core.ErrInvalidTxIn {
			return c.JSONError(ErrInvalidTxIn)
		} else if err == core.ErrTxOutAlreadySpent {
			return c.JSONError(ErrTxOutAlreadySpent)
		} else if err == core.ErrInvalidTxSignature {
			return c.JSONError(ErrInvalidTxSignature)
		} else if err == core.ErrInvalidTxBalance {
			return c.JSONError(ErrInvalidTxBalance)
		}
		return c.JSONError(err)
	}

	// COMMIT TX
	if err := txn.Commit(nil); err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(nil)
}

func (p *Plasma) GetTxHandler(c *Context) error {
	txPos, err := c.GetTxPositionFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	// BEGIN RO TXN
	txn := p.db.NewTransaction(false)
	defer txn.Discard()

	tx, err := p.childChain.GetTx(txn, txPos)
	if err != nil {
		if err == core.ErrTxNotFound {
			return c.JSONError(ErrTxNotFound)
		}
		return c.JSONError(err)
	}

	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]string{
		"tx": utils.EncodeToHex(txBytes),
	})
}

func (p *Plasma) GetTxProofHandler(c *Context) error {
	txPos, err := c.GetTxPositionFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	// BEGIN RO TXN
	txn := p.db.NewTransaction(false)
	defer txn.Discard()

	txProofBytes, err := p.childChain.GetTxProof(txn, txPos)
	if err != nil {
		if err == core.ErrTxNotFound {
			return c.JSONError(ErrTxNotFound)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]string{
		"proof": utils.EncodeToHex(txProofBytes),
	})
}
