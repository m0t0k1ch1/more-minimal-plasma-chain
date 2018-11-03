package app

func (p *Plasma) PingHandler(c *Context) error {
	return c.JSONSuccess(nil)
}
