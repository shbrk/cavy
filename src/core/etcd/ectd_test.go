package etcd

import (
	"fmt"
	"testing"
)

var cli *Client

func TestMain(m *testing.M) {
	client, err := NewClient(&Config{Endpoints: []string{"127.0.0.1:2379"}})
	if err != nil {
		fmt.Println(err)
		return
	}
	cli = client
	go func() {
		cli.Run()
	}()
	m.Run()
}

func TestSyncGet(t *testing.T) {
	exist, val, err := cli.SyncGet("/DEV/gate")
	if err != nil {
		t.Fail()
	}
	fmt.Println(exist, val)

	keys, values, err := cli.SyncGetWithPrefix("/a")
	if err == nil {
		for i := range keys {
			fmt.Println(keys[i], values[i])
		}
	} else {
		t.Fail()
	}
}

func TestSyncSet(t *testing.T) {
	err := cli.SyncSet("/a", "123")
	if err != nil {
		t.Fail()
	}
}

func TestAsyncGet(t *testing.T) {
	cli.AsyncGet("/a", func(err error, exist bool, key string, value string) {
		if err != nil {
			t.Fail()
		}
	})
	cli.AsyncGetWithPrefix("/a", func(err error, key string, keys []string, values []string) {
		if err != nil {
			t.Fail()
		}
	})
}

func TestAsyncSet(t *testing.T) {
	cli.AsyncSet("/watch", "456", func(err error) {
		if err != nil {
			t.Fail()
		}
	})
}

func TestKeepAlive(t *testing.T) {
	cli.KeepAlive("/alive", "abc", 5, func(err error, key string) {
		if err != nil {
			t.Fail()
		}
	})
}
func TestWatch(t *testing.T) {
	cli.Watch("/watch", true, func(err error, eventType EventType, key string, value string) {
		if err != nil {
			t.Fail()
		}
		fmt.Println("watch", err, eventType, key, value)
	})
}
