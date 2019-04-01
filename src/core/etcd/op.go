package etcd

import (
	"context"
	"go.etcd.io/etcd/clientv3"
)

type Op interface {
	Exec(c *Client)
}

type AsyncGetOp struct {
	Key  string
	Func func(err error, exist bool, key string, value string)
}

func (a *AsyncGetOp) Exec(c *Client) {
	ctx, cancel := context.WithTimeout(context.TODO(), c.timeout)
	defer cancel()
	resp, err := c.client.Get(ctx, a.Key)
	var event = &AsyncGetEvent{
		Key:  a.Key,
		Err:  err,
		Func: a.Func,
	}
	if err == nil {
		if resp.Kvs == nil || len(resp.Kvs) == 0 {
			event.Exist = false
		} else {
			event.Exist = true
			event.Value = string(resp.Kvs[0].Value)
		}
	}
	c.ChanOut <- event
}

type AsyncGetWithPrefixOp struct {
	Key  string
	Func func(err error, key string, keys []string, values []string)
}

func (a *AsyncGetWithPrefixOp) Exec(c *Client) {
	ctx, cancel := context.WithTimeout(context.TODO(), c.timeout)
	defer cancel()
	resp, err := c.client.Get(ctx, a.Key, clientv3.WithPrefix())
	keys := make([]string, 0)
	values := make([]string, 0)
	if err == nil {
		for _, val := range resp.Kvs {
			keys = append(keys, string(val.Key))
			values = append(keys, string(val.Value))
		}
	}
	var event = &AsyncGetWithPrefixEvent{
		Key:    a.Key,
		Err:    err,
		Keys:   keys,
		Values: values,
		Func:   a.Func,
	}
	c.ChanOut <- event
}

type AsyncSetOp struct {
	Key   string
	Value string
	Func  func(err error)
}

func (a *AsyncSetOp) Exec(c *Client) {
	ctx, cancel := context.WithTimeout(context.TODO(), c.timeout)
	defer cancel()
	_, err := c.client.Put(ctx, a.Key, a.Value)
	var event = &AsyncSetEvent{
		Key:  a.Key,
		Err:  err,
		Func: a.Func,
	}
	c.ChanOut <- event
}
