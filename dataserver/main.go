package main

import (
	"ESFS2.0/message/protos"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type dataServer struct {
	protos.UnimplementedAuthenticationServer
	protos.UnimplementedFileHandleServer
}

func main() {
	port := 8927
	host := "0.0.0.0"
	addr := fmt.Sprintf("%s:%d", host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[grpc] failed to listen: %v", err.Error())
	}
	defer lis.Close()

	s := grpc.NewServer()
	protos.RegisterAuthenticationServer(s, &dataServer{})
	protos.RegisterFileHandleServer(s, &dataServer{})

	go fileSocketServer()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("[grpc] failed to serve: %v", err)
	}

}
