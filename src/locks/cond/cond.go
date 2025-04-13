package cond

import (
	"log"
	"sync"
	"time"
)

func array(ii int) (v []int) {
	for i := 0; i < ii; i++ {
		v = append(v, i)
	}
	return
}

func F() {
	//S := array(11)
	c := sync.NewCond(&sync.Mutex{})

	//ch := make(chan int)

	for i, v := range array(2) {
		go func() {
			log.Printf("locking %d %d \n", i, v)
			c.L.Lock()
			log.Printf("locked %d %d \n", i, v)
			c.Wait()
			log.Printf("after cond.signal received. unlocking.\n")
			c.L.Unlock()
		}()
	}

	time.Sleep(time.Second * 2)

	for range 2 {
		<-time.Tick(time.Second)
		c.Signal()
	}
	time.Sleep(time.Second * 12)

	log.Printf("exit")

}
