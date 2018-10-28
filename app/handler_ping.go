package app

func (cc *ChildChain) PingHandler(c *Context) error {
	return c.JSONSuccess(nil)
}
