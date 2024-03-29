package message

import (
	"os"
	"time"
)

type OpType int8

const (
	FILE_UPLOAD   OpType = 1
	FILE_DOWNLOAD OpType = 2
)

type FileInfo struct {
	Name    string
	Mode    os.FileMode
	Size    int64
	ModTime time.Time
	Owner   string
}

type FileWithOwner struct {
	Filename string
	Owner    string
}

type FileSocketMessage struct {
	UserName  string
	FileName  []string
	Type      OpType
	AddOnInfo []FileWithOwner
}

type SignalOver struct {
}
