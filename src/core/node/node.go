package node

import (
	"context"
	"sync"
)

type Base struct {
	Ctx context.Context
	Wg  *sync.WaitGroup
}

type Node interface {
	Name() string
	Init()
	Run()
}
