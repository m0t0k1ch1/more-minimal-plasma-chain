package app

func (cc *ChildChain) PostTxHandler(c *Context) error {
	c.Request().ParseForm()

	tx, err := c.GetTxFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	if err := cc.mempool.Add(tx); err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(tx)
}
