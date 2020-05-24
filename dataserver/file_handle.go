package main

import (
	"ESFS2.0/dataserver/common"
	"ESFS2.0/message"
	"ESFS2.0/message/protos"
	"context"
	"encoding/json"
	"fmt"
	_ "fmt"
	"log"
	"net"
	"os"
	"path"
)

func fileSocketServer() {
	port := 8959
	host := "0.0.0.0"
	addr := fmt.Sprintf("%s:%d", host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[socket] failed to listen: %v", err.Error())
	}
	defer lis.Close()

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("建立socket连接失败 %v", err.Error())
			continue
		}

		msg := &message.FileSocketMessage{}
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		err = json.Unmarshal(buf[:n], msg)
		if err != nil {
			fmt.Printf("反序列化失败 %v", err.Error())
			continue
		}

		fmt.Println(msg.UserName, msg.FileName, msg.Type)

		//go func() {
		//
		//}()
	}

}

/**
@author js
*/
func (s *dataServer) UploadPrepare(ctx context.Context, req *protos.UploadPrepareRequest) (*protos.UploadPrepareResponse, error) {
	//反序列化，获取文件信息
	fileInfo := &message.FileInfo{}
	err := json.Unmarshal(req.FileInfo, fileInfo)
	if err != nil {
		log.Printf("反序列化失败 %v", err.Error())
		return &protos.UploadPrepareResponse{
			Ok:           false,
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
		}, err
	}

	//创建指定文件
	//Create函数若文件已存在则会截断，不存在则新建
	file, err := os.Create(path.Join(common.BaseDir, "dataserver", "data", req.Username, fileInfo.Name))
	if err != nil {
		log.Printf("创建文件失败 %v", err.Error())
		return &protos.UploadPrepareResponse{
			Ok:           false,
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
		}, err
	}
	defer file.Close()

	fmt.Println(fileInfo.Name, fileInfo.Size, fileInfo.Mode, fileInfo.ModTime)
	return &protos.UploadPrepareResponse{
		Ok:           true,
		ErrorMessage: protos.ErrorMessage_OK,
	}, nil
}

func (s *dataServer) ListFiles(ctx context.Context, req *protos.ListFilesRequest) (*protos.ListFilesResponse, error) {
	return nil, nil
}
