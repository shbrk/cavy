package etcd

import (
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type EventType = mvccpb.Event_EventType

const (
	EventType_PUT    = mvccpb.PUT
	EventType_DELETE = mvccpb.DELETE
)

type Event interface {
	HandleEvent()
}

//监视事件
type WatchEvent struct {
	Type  EventType
	Err   error
	Key   string
	Value string
	Func  func(err error, eventType EventType, key string, value string)
}

func (w *WatchEvent) HandleEvent() {
	if w.Func != nil {
		w.Func(w.Err, w.Type, w.Key, w.Value)
	}
}

//异步Set事件
type AsyncSetEvent struct {
	Err  error
	Key  string
	Func func(err error)
}

func (a *AsyncSetEvent) HandleEvent() {
	if a.Func != nil {
		a.Func(a.Err)
	}
}

//异步Get事件
type AsyncGetEvent struct {
	Err   error
	Exist bool
	Key   string
	Value string
	Func  func(err error, exist bool, key string, value string)
}

func (a *AsyncGetEvent) HandleEvent() {
	if a.Func != nil {
		a.Func(a.Err, a.Exist, a.Key, a.Value)
	}
}

//异步Get事件
type AsyncGetWithPrefixEvent struct {
	Err    error
	Key    string
	Keys   []string
	Values []string
	Func   func(err error, key string, keys []string, values []string)
}

func (a *AsyncGetWithPrefixEvent) HandleEvent() {
	if a.Func != nil {
		a.Func(a.Err, a.Key, a.Keys, a.Values)
	}
}

//异步接收keepalive结果
type KeepAliveEvent struct {
	Err  error
	Key  string
	Func func(err error, key string)
}

func (k *KeepAliveEvent) HandleEvent() {
	if k.Func != nil {
		k.Func(k.Err, k.Key)
	}
}
