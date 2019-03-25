package share

import (
	"container/heap"
	"time"
)

// TimerQueueSize TODO
const TimerQueueSize = 1024

// TimeOuter TODO
type TimeOuter interface {
	TimeOut(int64)
}

type Timer struct {
	onTimer         func(int64, interface{})
	CustomParameter interface{}
	index           int   // The index of the item in the heap.
	delay           int64 // The value of the item; arbitrary.
	point           int64 // The point of the item in the queue.
	interval        int
}

// A TimerHeap implements heap.Interface and holds Items.
type TimerHeap []*Timer

func (t TimerHeap) Len() int { return len(t) }

func (t TimerHeap) Less(i, j int) bool {
	return t[i].point < t[j].point
}

func (t TimerHeap) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
	t[i].index = i
	t[j].index = j
}

func (t *TimerHeap) Push(x interface{}) {
	n := len(*t)
	item := x.(*Timer)
	item.index = n
	*t = append(*t, item)
}

func (t *TimerHeap) Pop() interface{} {
	old := *t
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*t = old[0 : n-1]
	return item
}

// TODO 把timer 堆管理放在单独的线程管理
//精度取决于Run函数调用精度
type TimerManager struct {
	timeHeap   TimerHeap
	readyQueue []*Timer
}

func NewTimerManager() *TimerManager {
	return &TimerManager{timeHeap: make([]*Timer, 0, TimerQueueSize), readyQueue: make([]*Timer, 0, TimerQueueSize)}
}

// AddTimer 回调函数，时间间隔，循环次数，0表示没有次数限制，customParam回调函数传入的函数
func (m *TimerManager) AddTimer(onTimer func(now int64, param interface{}), delay time.Duration, interval int, param interface{}) *Timer {
	delayNum := delay.Nanoseconds()
	point := time.Now().UnixNano() + delayNum
	timer := &Timer{onTimer: onTimer, interval: interval, point: point, delay: delayNum, CustomParameter: param}
	heap.Push(&m.timeHeap, timer)
	return timer
}

func (m *TimerManager) RemoveTimer(t *Timer) {
	if t.index < 0 {
		return
	}
	heap.Remove(&m.timeHeap, t.index)
}

// now  unix nano
func (m *TimerManager) Run(now int64, limit int) {
	for len(m.timeHeap) > 0 {
		if m.timeHeap[0].point > now {
			break
		}
		timer := heap.Pop(&m.timeHeap).(*Timer)
		m.readyQueue = append(m.readyQueue, timer)
		timer.interval -= 1
		if timer.interval > 0 {
			timer.point = now + timer.delay
			heap.Push(&m.timeHeap, timer)
		}
		if limit > 0 && len(m.readyQueue) >= limit {
			break
		}
	}
	if len(m.readyQueue) == 0 {
		return
	}
	for i := range m.readyQueue {
		m.readyQueue[i].onTimer(now, m.readyQueue[i].CustomParameter)
	}
	m.readyQueue = m.readyQueue[:0] // 清空队列，TODO timer对象没有释放
}
