package concurrency

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func Fire(gmp int) {
	//fmt.Printf("runtime.NumCPU(): %d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(gmp)

	intchan := make(chan int)

	sg := sync.WaitGroup{}
	for i := range 100 {
		sg.Add(1)
		go func(i int) {
			defer sg.Done()
			//fmt.Printf("sending %d \n", i)
			//<-time.Tick(time.Duration(10-i) * time.Second)
			intchan <- i
		}(i)
	}

	go func() {
		sg.Wait()
		close(intchan)
	}()

	s := sync.WaitGroup{}

	for i := range 10 {
		sg.Add(1)
		go func(i int) {
			defer s.Done()
			for range intchan {
				//fmt.Printf("recv %d :val %d\n", i, v)
			}
		}(i)
	}

	s.Wait()
	//
	//exitSignal := make(chan os.Signal, 1)
	//signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	//
	//<-exitSignal
}

func PrintNumbers(name string, i int) {
	a := atomic.Int32{}
	a.Load()
	for ; i > 0; i-- {
		fmt.Printf("%s %d\n", name, i)
		runtime.Gosched()
	}
}
