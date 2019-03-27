package etcd

import (
	"context"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"time"
)

const IN_CHANNEL_SIZE = 1024
const OUT_CHANNEL_SIZE = 1024
const DEFAULT_TIME_OUT = 10 * time.Second

type Config struct {
	Endpoints []string //地址
	Timeout   int64    //超时时间 单位秒
	Username  string   //用户名
	Password  string   //密码
}

type Client struct {
	client  *clientv3.Client
	timeout time.Duration
	chanIn  chan Op
	chanOut chan Event
}


func NewClient(cfg *Config) (*Client, error) {
	var client = &Client{
		timeout: DEFAULT_TIME_OUT,
		chanIn:  make(chan Op, IN_CHANNEL_SIZE),
		chanOut: make(chan Event, OUT_CHANNEL_SIZE),
	}
	if cfg.Timeout != 0 {
		client.timeout = time.Duration(cfg.Timeout) * time.Second
	}
	var conf = clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: client.timeout,
		Username:    cfg.Username,
		Password:    cfg.Password,
	}
	cli, err := clientv3.New(conf)
	if err != nil {
		return nil, err
	}
	client.client = cli
	client.initThread()
	return client, nil
}

//维持一个key活跃 onErrCB可以传入在keep过程中出现的错误
func (c *Client) KeepAlive(key string, value string, timeout int64, onErrorCallback func(err error, key string)) {
	if timeout == 0 {
		timeout = int64(c.timeout / time.Second)
	}
	leaseResp, err := c.client.Grant(context.TODO(), timeout)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.TODO(), c.timeout)
	defer cancel()
	_, err = c.client.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		if onErrorCallback != nil {
			onErrorCallback(err, key)
		}
		return
	}
	ch, err := c.client.KeepAlive(context.TODO(), leaseResp.ID)
	//加入到监控
	if err != nil {
		if onErrorCallback != nil {
			onErrorCallback(err, key)
		}
		return
	}
	go func() {
		for {
			_, ok := <-ch
			if !ok {
				event := &KeepAliveEvent{Err: errors.New("keep alive closed"), Key: key, Func: onErrorCallback}
				c.chanOut <- event
			}
		}
	}()
	return
}

//监控一个key的变动 withPrefix打开后监听的是文件夹 onEventCB是事件触发后执行的回调函数
func (c *Client) Watch(key string, withPrefix bool, onEventCallback func(err error, eventType EventType, key string, value string)) {
	go func() {
		var wc clientv3.WatchChan
		if withPrefix {
			wc = c.client.Watch(context.TODO(), key, clientv3.WithPrefix(), clientv3.WithPrevKV())
		} else {
			wc = c.client.Watch(context.TODO(), key)
		}
		for {
			ws, ok := <-wc
			if !ok || ws.Canceled {
				event := &WatchEvent{Err: errors.New("watch closed"), Key: key, Func: onEventCallback}
				c.chanOut <- event
				return
			} else {
				for _, e := range ws.Events {
					event := &WatchEvent{Type: EventType(e.Type), Key: string(e.Kv.Key),Value:string(e.Kv.Value),Func: onEventCallback}
					c.chanOut <- event
				}
			}
		}
	}()
}

//同步设置一个key
func (c *Client) SyncSet(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), c.timeout)
	_, err := c.client.Put(ctx, key, value)
	cancel()
	return err
}

//同步设置一个key
func (c *Client) SyncGet(key string) (exist bool, value string, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), c.timeout)
	resp, err := c.client.Get(ctx, key)
	cancel()
	if err == nil {
		if resp.Kvs == nil || len(resp.Kvs) == 0 {
			exist = false
		} else {
			exist = true
			value = string(resp.Kvs[0].Value)
		}
	}
	return
}

//获取一个路径的所有key
func (c *Client) SyncGetWithPrefix(key string) (keys []string, values []string, err error) {
	resp, err := c.client.Get(context.TODO(), key, clientv3.WithPrefix())
	keys = make([]string, 0)
	values = make([]string, 0)
	if err == nil {
		for _, val := range resp.Kvs {
			keys = append(keys, string(val.Key))
			values = append(values, string(val.Value))
		}
	}
	return
}

//异步获取一个key
func (c *Client) AsyncGet(key string, onValueCallback func(err error, exist bool, key string, value string)) {
	c.produce(&AsyncGetOp{key, onValueCallback})
}

func (c *Client) AsyncGetWithPrefix(key string, onValueCallback func(err error, key string, keys []string, values []string)) {
	c.produce(&AsyncGetWithPrefixOp{key, onValueCallback})
}

//异步设置一个key
func (c *Client) AsyncSet(key string, value string, onErrorCallback func(err error)) {
	c.produce(&AsyncSetOp{key, value, onErrorCallback})
}
