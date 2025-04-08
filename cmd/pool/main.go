package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	p := sync.Pool{
		New: func() func() interface{} {
			a := 1
			return func() interface{} {
				defer func() { a++ }()
				return a
			}
		}(),
	}

	fmt.Println(p.Get())
	fmt.Println(p.Get())
	fmt.Println(p.Get())
	p.Put(222)

	runtime.GC()
	<-time.Tick(time.Second)
	fmt.Println(p.Get())
}
