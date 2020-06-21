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
	"time"
)

type MyRegisterWindow struct {
	*walk.MainWindow
	username *walk.LineEdit
	pwd      *walk.LineEdit
	repwd    *walk.LineEdit
	dft2pwd  *walk.LineEdit
	certPath *walk.LineEdit
}

func register(_username, _pwd, _repwd, _dft2pwd, _certPath string) {
	username, pwd, repwd, dft2pwd, certPath := _username, _pwd, _repwd, _dft2pwd, _certPath
	if username == "" || pwd == "" || repwd == "" || dft2pwd == "" || certPath == "" {
		common.ShowMsgBox("提示", "请完整填写信息")
		return
	}
	if pwd != repwd {
		common.ShowMsgBox("提示", "两次密码不一致")
		return
	}

	c, conn, err := common.GetAuthenticationClient()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cert, err := utils.ReadCertFromFile(certPath)
	if err != nil {
		common.ShowMsgBox("提示", "解析证书文件失败")
		return
	}
	if cert.Info.Username != username {
		common.ShowMsgBox("提示", "证书所有者与用户名不一致")
		return
	}
	serializedCert, err := json.Marshal(cert)
	if err != nil {
		common.ShowMsgBox("提示", "程序出错")
		return
	}
	request := &protos.RegisterRequest{
		Username:         username,
		Password:         pwd,
		DefaultSecondKey: dft2pwd,
		CertData:         serializedCert,
	}

	response, err := c.Register(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	switch response.ErrorMessage {
	case protos.ErrorMessage_OK:
		common.ShowMsgBox("提示", "注册成功")
	case protos.ErrorMessage_USER_ALREADY_EXISTS:
		common.ShowMsgBox("提示", "该用户名已存在")
	case protos.ErrorMessage_SERVER_ERROR:
		common.ShowMsgBox("提示", "服务器错误")
	}
}

func CreateRegisterWindow() {
	mw := &MyRegisterWindow{}
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "账号注册",
		Size:     Size{450, 100},
		Layout:   VBox{},
		Children: []Widget{
			GroupBox{
				//MaxSize: Size{500, 500},
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
					Label{Text: "密 码"},
					LineEdit{
						AssignTo:     &mw.pwd,
						PasswordMode: true,
					},
				},
			},
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "确认密码"},
					LineEdit{
						AssignTo:     &mw.repwd,
						PasswordMode: true,
					},
				},
			},
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "默认二级密码"},
					LineEdit{
						AssignTo:     &mw.dft2pwd,
						PasswordMode: true,
					},
				},
			},
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &mw.certPath,
						ReadOnly: true,
					},
					PushButton{
						Text: "选择证书文件",
						OnClicked: func() {
							mw.chooseCert() //TODO:选择证书
						},
					},
				},
			},
			PushButton{
				Text: "注册",
				OnClicked: func() {
					register(mw.username.Text(), mw.pwd.Text(), mw.repwd.Text(), mw.dft2pwd.Text(), mw.certPath.Text())
				},
			},
			PushButton{
				Text: "转到keyserver的注册界面",
				OnClicked: func() {
					CreateKeyServerWindow()
				},
			},
		},
	}.Create()); err != nil {
		fmt.Println("Error!")
	}
}

func (mw *MyRegisterWindow) chooseCert() {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择证书文件"
	if ok, err := dlg.ShowOpen(mw); err != nil {
		return
	} else if !ok {
		return
	} else {
		mw.certPath.SetText(dlg.FilePath)
	}
}
