package etcd

func (c *Client) initThread() {
	go func() {
		for {
			op := <-c.chanIn
			op.Exec(c)
		}
	}()
}

//主线程需要循环调用的消费函数
func (c *Client) Run() {
	select {
	case event, ok := <-c.chanOut:
		if ok {
			event.HandleEvent()
		}
	default:
		return
	}
}
func (c *Client) produce(op Op) {
	c.chanIn <- op
}
