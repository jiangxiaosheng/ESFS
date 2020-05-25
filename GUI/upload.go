package main

import (
	"ESFS2.0/client"
	"ESFS2.0/message"
	"ESFS2.0/message/protos"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MyMainWindow struct {
	*walk.MainWindow
	selectedfile *walk.LineEdit
	message      *walk.ListBox
	wv           *walk.WebView
	textEdit     *walk.TextEdit
}

var _useRadioButton, _notuseRadioButton *walk.RadioButton

func OpenWindow() {
	mw := &MyMainWindow{}
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "文件上传",

		MinSize: Size{300, 400},
		Layout:  VBox{},
		Children: []Widget{
			GroupBox{
				MaxSize: Size{500, 500},
				Layout:  HBox{},
				Children: []Widget{
					PushButton{ //选择文件
						Text:      "打开文件",
						OnClicked: mw.selectFile, //点击事件
					},
					Label{Text: "选中的文件 "},
					LineEdit{
						AssignTo: &mw.selectedfile, //选中的文件
					},
					PushButton{
						Text: "上传",
						OnClicked: func() {
							upload(mw)
						},
					},
				},
			},
			GroupBox{
				MaxSize: Size{500, 500},
				Layout:  HBox{},
				Children: []Widget{
					RadioButton{
						AssignTo: &_useRadioButton,
						Text:     "使用默认二级密码",
						OnClicked: func() {
							mw.selectedfile.SetReadOnly(true)
						},
					},
					RadioButton{
						AssignTo: &_notuseRadioButton,
						Text:     "使用新二级密码",
						OnClicked: func() {
							mw.selectedfile.SetReadOnly(false)
						},
					},
					LineEdit{
						AssignTo: &mw.selectedfile, //新的二级密码
						ReadOnly: true,
					},
				},
			},
			TextEdit{
				//MaxSize: Size{300, 200},
				AssignTo: &mw.textEdit,
				ReadOnly: true,
				Text:     "Drop files here, from windows explorer...",
			},
		},
	}.Create()); err != nil {
		fmt.Printf("Run err: %+v\n", err)
	}
	mw.textEdit.DropFiles().Attach(func(files []string) {
		mw.textEdit.SetText(strings.Join(files, "\r\n"))
	})
	mw.Run()
}

func upload(mw *MyMainWindow) {
	path := mw.selectedfile.Text()
	file, err := os.Open(path)
	if err != nil {
		log.Printf("打开文件失败 %v", err.Error())
		ShowMsgBox("提示", "打开文件失败")
		return
	}

	//异步上传
	go func() {
		//通过grpc向dataserver发送准备传输请求，这一步获得用户的默认二级密码
		c, grpcConn, err := client.GetFileHandleClient()
		if err != nil {
			log.Printf("建立grpc客户端失败 %v", err.Error())
			ShowMsgBox("提示", "服务器错误")
			return
		}
		defer grpcConn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("读取文件信息失败 %v", err.Error())
			ShowMsgBox("提示", "读取文件信息失败")
			return
		}

		buf, err := json.Marshal(fileInfo)
		if err != nil {
			log.Printf("序列化文件信息失败 %v", err.Error())
			ShowMsgBox("提示", "读取文件信息失败")
			return
		}

		prepareRequest := &protos.UploadPrepareRequest{
			Username: CurrentUser,
			FileInfo: buf,
		}
		response, err := c.UploadPrepare(ctx, prepareRequest)
		if err != nil {
			log.Printf("服务器错误 %v", err.Error())
			ShowMsgBox("提示", "服务器错误")
			return
		}

		switch response.ErrorMessage {
		case protos.ErrorMessage_SERVER_ERROR:
			log.Printf("服务器错误 %v", err.Error())
			ShowMsgBox("提示", "服务器错误")
			return
		case protos.ErrorMessage_OK:
			//通过socket发送实际的文件数据
			addr := "0.0.0.0:8959"
			socketConn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Println(err)
			}
			defer socketConn.Close()

			msg := message.FileSocketMessage{
				UserName: CurrentUser,
				FileName: filepath.Base(file.Name()),
				Type:     message.FILE_UPLOAD,
			}

			serializedData, err := json.Marshal(msg)
			_, err = socketConn.Write(serializedData)
			if err != nil {
				fmt.Printf("socket写入数据失败 %v", err.Error())
				return
			}

			buffer := make([]byte, 2048)
			for {
				n, err := file.Read(buffer)
				if err == io.EOF {
					break
				}
				_, err = socketConn.Write(buffer[:n])
				if err != nil {
					fmt.Println(err)
					break
				}
			}

			ShowMsgBox("提示", "上传成功")
		}
	}()
}

func (mw *MyMainWindow) selectFile() {
	allowType := viper.GetStringSlice("common.server.allowtype")
	fmt.Printf("allowType: %+v\n", allowType)

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	//dlg.Filter = "可上传jpg (*.jpg)|*.jpg|可上传png (*.png)|*.png|可上传gif (*.gif)|*.gif|所有文件 (*.*)|*.*"
	//判断可允许上传的文件
	filter := []string{}
	filterstring := ""
	for _, v := range allowType {
		if v != "*" {
			filterstring = "可上传" + v + " (*." + v + ")|*." + v
			filter = append(filter, filterstring)
		} else {
			filterstring = "所有文件" + v + " (*." + v + ")|*." + v
			filter = append(filter, filterstring)
		}
	}
	dlg.Filter = strings.Join(filter, "|") //切片转换字符串
	fmt.Printf("dlg.Filter: %+v\n", dlg.Filter)

	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.selectedfile.SetText("") //通过重定向变量设置TextEdit的Text
		return
	} else if !ok {
		mw.selectedfile.SetText("") //通过重定向变量设置TextEdit的Text
		return
	}
	mw.selectedfile.SetText(dlg.FilePath) //通过重定向变量设置TextEdit的Text

	//************************文件path**************************************************
	fmt.Printf("dlg.FilePath: %+v\n", dlg.FilePath)
	//********************************************************************************
}
