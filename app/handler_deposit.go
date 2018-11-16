package app

func (p *Plasma) PostDepositHandler(c *Context) error {
	c.Request().ParseForm()

	ownerAddr, err := c.GetOwnerAddressFromForm()
	if err != nil {
		return c.JSONError(err)
	}
	amount, err := c.GetAmountFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	// BEGIN TXN
	txn := p.db.NewTransaction(true)
	defer txn.Discard()

	newBlkNum, err := p.childChain.AddDepositBlock(txn, ownerAddr, amount, p.operator)
	if err != nil {
		return c.JSONError(err)
	}

	// COMMIT TXN
	if err := txn.Commit(nil); err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]uint64{
		"blknum": newBlkNum,
	})
}
