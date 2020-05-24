package client

import (
	"ESFS2.0/message/protos"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func GetAuthenticationClient() (protos.AuthenticationClient, *grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8927)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("无法建立grpc连接 %v", err.Error())
		return nil, conn, err
	}

	c := protos.NewAuthenticationClient(conn)
	return c, conn, nil
}

func GetFileHandleClient() (protos.FileHandleClient, *grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8927)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("无法建立grpc连接 %v", err.Error())
		return nil, conn, err
	}

	c := protos.NewFileHandleClient(conn)
	return c, conn, nil
}
