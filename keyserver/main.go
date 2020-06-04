package main

import (
	"ESFS2.0/message/protos"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type keyServer struct {
	protos.UnimplementedCAServer
}

func main() {
	port := 9015
	host := "0.0.0.0"
	addr := fmt.Sprintf("%s:%d", host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[grpc] failed to listen: %v", err.Error())
	}
	defer lis.Close()

	s := grpc.NewServer()
	protos.RegisterCAServer(s, &keyServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("[grpc] failed to serve: %v", err)
	}
}
