package main

func (cc *ChildChain) PingHandler(c *Context) error {
	return c.JSONSuccess(nil)
}
