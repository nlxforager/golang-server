package main

import (
	"runtime"
	"time"

	"golang-server/src/concurrency"
)

func main() {
	runtime.GOMAXPROCS(4)
	//go func() {
	go concurrency.PrintNumbers("g1", 50)
	go concurrency.PrintNumbers("g2", 50)
	//}()/

	time.Sleep(200 * time.Second)
}
