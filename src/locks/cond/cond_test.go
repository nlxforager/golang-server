package cond

import (
	"sync"
	"testing"
)

func TestHello(t *testing.T) {

	S := []int{1, 2, 3, 4}

	c := sync.NewCond(&sync.Mutex{})
	for i, s := range S {
		c.L.Lock()
		_, _ = i, s
		c.Signal()
	}
}
