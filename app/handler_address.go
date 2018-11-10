package app

import "github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"

func (p *Plasma) GetAddressUTXOsHandler(c *Context) error {
	addr, err := c.GetAddressFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	// BEGIN RO TXN
	txn := p.db.NewTransaction(false)
	defer txn.Discard()

	utxoPoses, err := p.childChain.GetUTXOPositions(txn, addr)
	if err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string][]types.Position{
		"utxos": utxoPoses,
	})
}
