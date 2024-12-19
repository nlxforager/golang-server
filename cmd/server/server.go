package main

import (
	"context"
	"log"
)

func main() {
	ctx := context.Background()
	log.Println("starting server")
	select {
	default:
		_ = ctx
	}

	log.Println("stopping server")
}
