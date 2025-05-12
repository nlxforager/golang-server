package main

import (
	"context"
	"fmt"
	"time"
)

func foo(ctx context.Context, i int) {
	fmt.Printf("foo %d\n", i)
	if i == 4 {

		foo(context.Background(), i+1)
		return
	}
	if i == 6 {
		return
	}
	go foo(ctx, i+1)
	<-ctx.Done()
	fmt.Printf("foo %d cancelled\n", i)
}
func main() {
	//
	//
	go bar(context.Background())
	time.Sleep(time.Second * 300)
	select {}
}

func bar(ctx context.Context) {
	_ctx, _ := context.WithDeadline(ctx, time.Now().Add(time.Second*2))
	go foo(_ctx, 0)

	<-ctx.Done()
	fmt.Printf("%d bar cancelled\n", 0)
}
