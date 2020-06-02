package GUI

import (
	clicommon "ESFS2.0/client/common"
	"ESFS2.0/message"
	"ESFS2.0/message/protos"
	"ESFS2.0/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
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
	secondpw     *walk.LineEdit
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
							mw.secondpw.SetReadOnly(true)
						},
					},
					RadioButton{
						AssignTo: &_notuseRadioButton,
						Text:     "使用新二级密码",
						OnClicked: func() {
							mw.secondpw.SetReadOnly(false)
						},
					},
					LineEdit{
						AssignTo: &mw.secondpw, //新的二级密码
						ReadOnly: true,
					},
				},
			},
			TextEdit{
				AssignTo: &mw.textEdit,
				ReadOnly: true,
				Text:     "拖拽文件到此处",
			},
		},
	}.Create()); err != nil {
		fmt.Printf("Run err: %+v\n", err)
	}
	mw.textEdit.DropFiles().Attach(func(files []string) {
		mw.textEdit.SetText(strings.Join(files, "\r\n"))
	})
	_useRadioButton.SetChecked(true)
	mw.Run()
}

func upload(mw *MyMainWindow) {
	if _notuseRadioButton.Checked() == true && mw.secondpw.Text() == "" {
		clicommon.ShowMsgBox("提示", "请填写二级密码")
		return
	}

	path := mw.selectedfile.Text()
	file, err := os.Open(path)
	if err != nil {
		log.Printf("打开文件失败 %v", err.Error())
		clicommon.ShowMsgBox("提示", "打开文件失败")
		return
	}

	go func() {
		//1.建立grpc连接，发送文件准备请求，同时得到默认二级密码
		c, conn, err := clicommon.GetFileHandleClient()
		if err != nil {
			fmt.Println(err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		stat, _ := file.Stat()
		fileInfo := message.FileInfo{
			Name:    stat.Name(),
			Mode:    stat.Mode(),
			Size:    stat.Size(),
			ModTime: stat.ModTime(),
		}
		serializedData, err := json.Marshal(fileInfo)
		if err != nil {
			log.Printf("序列化文件信息失败 %v", err.Error())
			clicommon.ShowMsgBox("提示", "服务器错误")
			return
		}

		request := &protos.UploadPrepareRequest{
			Username: CurrentUser,
			FileInfo: serializedData,
		}

		//上传二级密码
		if _useRadioButton.Checked() {
			request.SecondKey = ""
		} else {
			request.SecondKey = mw.secondpw.Text()
		}

		response, err := c.UploadPrepare(ctx, request)

		switch response.ErrorMessage {
		case protos.ErrorMessage_SERVER_ERROR:
			log.Printf("服务器错误 %v", err.Error())
			clicommon.ShowMsgBox("提示", "服务器错误")
			return
		case protos.ErrorMessage_OK:
			defaultSecondKey := response.DefaultSecondKey
			if _notuseRadioButton.Checked() {
				defaultSecondKey = mw.secondpw.Text()
			}
			log.Println(defaultSecondKey)

			//2.对文件进行AES加密，并生成签名信息
			privateKey := clicommon.GetUserPrivateKey()
			sessionKey, err := utils.GenerateSessionKeyWithSecondKey(defaultSecondKey, privateKey)
			if err != nil {
				log.Printf("服务器错误 %v", err.Error())
				clicommon.ShowMsgBox("提示", "服务器错误")
				return
			}
			encryptedData, err := utils.AESEncryptFileToBytes(path, sessionKey)
			if err != nil {
				log.Printf("加密文件失败 %v", err.Error())
				clicommon.ShowMsgBox("提示", "服务器错误")
				return
			}

			//3.建立socket连接，发送加密文件数据
			addr := "0.0.0.0:8959"
			socketConn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Println(err)
				clicommon.ShowMsgBox("提示", "服务器错误")
			}

			msg := message.FileSocketMessage{
				UserName: CurrentUser,
				FileName: []string{filepath.Base(file.Name())},
				Type:     message.FILE_UPLOAD,
			}

			serializedData, err = json.Marshal(msg)
			if err != nil {
				log.Printf("序列化失败 %v", err.Error())
				return
			}
			_, err = socketConn.Write(serializedData)
			if err != nil {
				log.Printf("socket写入数据失败 %v", err.Error())
				return
			}

			bytesReader := bytes.NewReader(encryptedData)
			buffer := make([]byte, 2048)
			socketConn.Read(buffer) //防止粘包
			for {
				n, err := bytesReader.Read(buffer) //从字节数组读数据
				if err == io.EOF {
					break
				}
				_, err = socketConn.Write(buffer[:n])
				if err != nil {
					fmt.Println(err)
					break
				}
			}
			if socketConn.Close() == nil {
				fmt.Println("close")
			} else {
				fmt.Println("fail")
			}

			//4.使用grpc服务传签名信息（因为签名数据不大，所以用grpc是没问题的）
			sigData, err := utils.SignatureFile(file, clicommon.GetUserPrivateKey())
			if err != nil {
				log.Printf("签名失败 %v", err.Error())
				clicommon.ShowMsgBox("提示", "服务器错误")
				return
			}

			dsRequest := &protos.UploadDSRequest{
				Username: CurrentUser,
				Filename: stat.Name(),
				DsData:   sigData,
			}
			dsResponse, err := c.UploadDS(ctx, dsRequest)
			switch dsResponse.ErrorMessage {
			case protos.ErrorMessage_SERVER_ERROR:
				clicommon.ShowMsgBox("提示", "服务器错误")
				return
			case protos.ErrorMessage_OK:
				clicommon.ShowMsgBox("提示", "上传成功")
				return
			}
		}
	}()
}

func (mw *MyMainWindow) selectFile() {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	//dlg.ShowSave(mw)
	filter := []string{}
	filterstring := "所有文件(*.*)"
	filter = append(filter, filterstring)
	dlg.Filter = strings.Join(filter, "|") //切片转换字符串
	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.selectedfile.SetText("") //通过重定向变量设置TextEdit的Text
		return
	} else if !ok {
		mw.selectedfile.SetText("") //通过重定向变量设置TextEdit的Text
		return
	}
	mw.selectedfile.SetText(dlg.FilePath) //通过重定向变量设置TextEdit的Text
}
