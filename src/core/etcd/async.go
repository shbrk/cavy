package etcd

func (c *Client) initThread() {
	go func() {
		for {
			op := <-c.chanIn
			op.Exec(c)
		}
	}()
}

func (c *Client) produce(op Op) {
	c.chanIn <- op
}
