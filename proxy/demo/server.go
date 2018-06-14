package main

import (
	"fmt"
	"net"
	"go-envoy-poc/proxy"
	"log"
	"google.golang.org/grpc/reflection"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	lis, err := net.Listen("tcp", ":9933")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := proxy.NewGrpcServer()
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


