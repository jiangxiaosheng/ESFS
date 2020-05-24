package tests

import (
	"ESFS2.0/client"
	"ESFS2.0/message"
	"ESFS2.0/message/protos"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestFileServer(t *testing.T) {

}

func TestFileClient(t *testing.T) {
	c, conn, err := client.GetFileHandleClient()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	file, _ := os.Open("auth_test.go")
	stat, _ := file.Stat()
	fileInfo := message.FileInfo{
		Name:    stat.Name(),
		Mode:    stat.Mode(),
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
	}
	serializedData, _ := json.Marshal(fileInfo)

	fmt.Println(stat.Name(), stat.Mode(), stat.Size(), stat.ModTime())

	request := &protos.UploadPrepareRequest{
		Username: "memeshe",
		FileInfo: serializedData,
	}

	response, err := c.UploadPrepare(ctx, request)
	if response != nil {
		fmt.Println(response.Ok, response.ErrorMessage)
	}
}
