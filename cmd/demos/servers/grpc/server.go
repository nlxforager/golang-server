package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"golang-server/cmd/servers/grpc/helloworldpb"
	"google.golang.org/grpc"
)

type greeterServer struct {
	helloworldpb.GreeterServer
}

var _ helloworldpb.GreeterServer = &greeterServer{}

func (g *greeterServer) SayHello(_ context.Context, r *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
	return &helloworldpb.HelloReply{
		Message: fmt.Sprintf("Hello %s im server", r.Name),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	defer lis.Close()

	server := grpc.NewServer()
	helloworldpb.RegisterGreeterServer(server, &greeterServer{})
	err = server.Serve(lis)
	if err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
