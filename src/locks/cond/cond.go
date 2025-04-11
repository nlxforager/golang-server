package cond

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func array(i int) (v []int) {
	for i > 0 {
		v = append(v, i)
		i -= 1
	}
	return
}

func F() {
	S := array(11)
	c := sync.NewCond(&sync.Mutex{})

	ch := make(chan int)
	for i := range 5 {
		//c.L.Lock()
		go func(i int) {
			//defer c.L.Unlock()
			for {
				log.Printf("[%d ] wait on ch \n", i)
				v := <-ch
				//log.Printf("[%d] v %d recv. waiting\n", i, v)
				c.Wait()
				//log.Printf("[%d] v %d recv. await finished\n", i, v)

				log.Printf("[%d] Got value = %d from channel\n", i, v)
			}
		}(i)
	}

	for i, s := range S {
		time.Sleep(200 * time.Millisecond)
		_, _ = i, s

		go func() {
			c.L.Lock()

			ch <- s
			fmt.Printf("SIG[%d] \n", i)
			c.Signal()
		}()

		//c.L.Unlock()
		//c.Broadcast()
		//c.Signal()
	}

	time.Sleep(time.Second * 1000)

}
