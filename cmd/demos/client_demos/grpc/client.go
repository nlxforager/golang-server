package main

import (
	"context"
	"log"
	"time"

	"golang-server/cmd/demos/servers/grpc/helloworldpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	gc, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := helloworldpb.NewGreeterClient(gc)

	for range time.NewTicker(1 * time.Second).C {
		v, err := client.SayHello(context.Background(), &helloworldpb.HelloRequest{
			Name: "John Dssoe",
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Println(v)
	}
}
