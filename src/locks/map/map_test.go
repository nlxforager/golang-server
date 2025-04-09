package _map

import (
	"sync"
	"testing"
)

func TestMutex(t *testing.T) {
	var m = make(map[string]int)
	l := sync.Mutex{}

	for i := 0; i < 100; i++ {
		go func(i int) {
			l.Lock()
			v, _ := m["123123123"]
			delete(m, "123123123")
			m["123123123"] = v + i
			l.Unlock()
		}(i)
	}
}

func TestRWMutex(t *testing.T) {
	//var m = make(map[string]int)

	//l := sync.RWMutex{}
	//for i := 0; i < 100; i++ {
	//	go func(i int) {
	//		l.RLock()
	//		v, _ := m["123123123"]
	//		delete(m, "123123123")
	//		m["123123123"] = v + i
	//		l.RUnlock()
	//	}(i)
	//}
}
