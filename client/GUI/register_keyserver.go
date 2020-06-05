package GUI

import (
	"ESFS2.0/client/common"
	"ESFS2.0/message/protos"
	"ESFS2.0/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"log"
	"path"
	"time"
)

type MyKeyserverRegisterWindow struct {
	*walk.MainWindow
	username    *walk.LineEdit
	selectpub   *walk.LineEdit
	generatecer *walk.LineEdit
}

func CreateKeyServerWindow() {
	mw := &MyKeyserverRegisterWindow{}
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Keyserver注册",
		Size:     Size{450, 100},
		Layout:   VBox{},
		Children: []Widget{
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "用户名"},
					LineEdit{
						AssignTo: &mw.username,
					},
				},
			},
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &mw.selectpub,
						ReadOnly: true,
					},
					PushButton{
						Text: "选择公钥文件",
						OnClicked: func() {
							mw.SelectPub() //TODO:选择证书
						},
					},
				},
			},
			PushButton{
				Text: "生成RSA密钥对",
				OnClicked: func() {
					mw.GenerateKey() //TODO:生成证书
				},
			},
			PushButton{
				Text: "注册",
				OnClicked: func() {
					mw.registerInCA() //TODO:keyserver注册
				},
			},
		},
	}.Create()); err != nil {
		fmt.Println("Error!")
	}
}

/**
@author yyx
*/
func (mw *MyKeyserverRegisterWindow) SelectPub() {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择公钥文件位置"
	dlg.Filter = "Certificate File(*.pem)|*.pem" //
	if ok, err := dlg.ShowOpen(mw); err != nil || !ok {
		mw.selectpub.SetText("")
		return
	} else {
		mw.selectpub.SetText(dlg.FilePath)
		return
	}
}

/**
@author yyx
*/
func (mw *MyKeyserverRegisterWindow) GenerateKey() {
	if mw.username.Text() == "" {
		common.ShowMsgBox("提示", "请输入用户名")
		return
	}

	dlg := new(walk.FileDialog)
	dlg.Title = "选择密钥对的保存位置"
	dlg.InitialDirPath = "Default"
	if ok, err := dlg.ShowBrowseFolder(mw); err != nil {
		log.Print(err)
		common.ShowMsgBox("提示", "程序错误")
	} else if ok == true {
		err = utils.GenerateRSAKey(1024, dlg.FilePath, mw.username.Text())
		if err != nil {
			log.Println(err)
			common.ShowMsgBox("提示", "生成密钥对出错")
			return
		}
		common.ShowMsgBox("提示", "生成密钥对成功")
		log.Println(dlg.FilePath)
	}
	return
}

func (mw *MyKeyserverRegisterWindow) registerInCA() {
	if mw.username.Text() == "" {
		common.ShowMsgBox("提示", "请输入用户名")
		return
	}

	if mw.selectpub.Text() == "" {
		common.ShowMsgBox("提示", "请选择公钥文件")
		return
	}

	dlg := new(walk.FileDialog)
	dlg.Title = "选择证书的保存位置"
	dlg.InitialDirPath = "Default"
	if ok, err := dlg.ShowBrowseFolder(mw); err != nil || !ok {
		return
	} else {
		pubKey := utils.GetPublicKeyFromFile(mw.selectpub.Text())

		if pubKey == nil {
			common.ShowMsgBox("提示", "解析公钥文件失败")
			return
		}

		serializedData, err := json.Marshal(pubKey)
		if err != nil {
			common.ShowMsgBox("提示", "程序错误")
			return
		}

		//grpc发送请求
		c, conn, err := common.GetCAClient()
		if err != nil {
			fmt.Println(err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		request := &protos.SetCertRequest{
			Username: mw.username.Text(),
			Content:  serializedData,
		}

		response, err := c.SetCert(ctx, request)
		switch response.ErrorMessage {
		case protos.ErrorMessage_USER_ALREADY_EXISTS:
			common.ShowMsgBox("提示", "用户已存在")
			return
		case protos.ErrorMessage_SERVER_ERROR:
			common.ShowMsgBox("提示", "服务器错误")
			return
		case protos.ErrorMessage_OK:
			err = ioutil.WriteFile(path.Join(dlg.FilePath, request.Username+"_cert.pem"), response.CertData, 0644)
			if err != nil {
				log.Println(err)
				common.ShowMsgBox("提示", "证书写入失败")
			} else {
				common.ShowMsgBox("提示", "KeyServer注册成功")
				_ = mw.Close()
			}
			return
		}
	}
}
