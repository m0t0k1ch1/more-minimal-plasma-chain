package app

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

func (p *Plasma) PostTxHandler(c *Context) error {
	c.Request().ParseForm()

	tx, err := c.GetTxFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	txPos, err := p.childChain.AddTxToMempool(tx)
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

	return c.JSONSuccess(map[string]types.Position{
		"pos": txPos,
	})
}

func (p *Plasma) GetTxHandler(c *Context) error {
	txPos, err := c.GetTxPositionFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	tx, err := p.childChain.GetTx(txPos)
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

	txProofBytes, err := p.childChain.GetTxProof(txPos)
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
