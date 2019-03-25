package share

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTimerManager(t *testing.T) {
	m := NewTimerManager()
	var i = 0
	timer := m.AddTimer(func(now int64, param interface{}) {
		fmt.Println(now, time.Now())
		i++
	}, time.Second, 5, nil)
	for {
		if i == 2 {
			m.RemoveTimer(timer)
			break
		}
		m.Run(time.Now().UnixNano(), 0)
	}
}
