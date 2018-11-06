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

	txHash, err := p.childChain.AddTxToMempool(tx)
	if err != nil {
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

	return c.JSONSuccess(map[string]interface{}{
		"txhash": utils.HashToHex(txHash),
	})
}

func (p *Plasma) GetTxHandler(c *Context) error {
	txHash, err := c.GetTxHashFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	tx, err := p.childChain.GetTx(txHash)
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

	return c.JSONSuccess(map[string]interface{}{
		"tx": utils.EncodeToHex(txBytes),
	})
}

func (p *Plasma) GetTxIndexHandler(c *Context) error {
	txHash, err := c.GetTxHashFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	blkNum, txIndex, err := p.childChain.GetTxIndex(txHash)
	if err != nil {
		if err == core.ErrTxNotFound {
			return c.JSONError(ErrTxNotFound)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]interface{}{
		"blknum":  blkNum,
		"txindex": txIndex,
	})
}

func (p *Plasma) GetTxProofHandler(c *Context) error {
	txHash, err := c.GetTxHashFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	txProofBytes, err := p.childChain.GetTxProof(txHash)
	if err != nil {
		if err == core.ErrTxNotFound {
			return c.JSONError(ErrTxNotFound)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]interface{}{
		"proof": utils.EncodeToHex(txProofBytes),
	})
}

func (p *Plasma) PutTxHandler(c *Context) error {
	c.Request().ParseForm()

	txHash, err := c.GetTxHashFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	iIndex, err := c.GetInputIndexFromForm()
	if err != nil {
		return c.JSONError(err)
	}
	confSig, err := c.GetConfirmationSignatureFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	if err := p.childChain.ConfirmTx(txHash, iIndex, confSig); err != nil {
		if err == core.ErrInvalidTxConfirmationSignature {
			return c.JSONError(ErrInvalidTxConfirmationSignature)
		} else if err == core.ErrTxInNotFound {
			return c.JSONError(ErrTxInNotFound)
		} else if err == core.ErrNullTxInConfirmation {
			return c.JSONError(ErrNullTxInConfirmation)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]interface{}{
		"txhash": utils.HashToHex(txHash),
	})
}
