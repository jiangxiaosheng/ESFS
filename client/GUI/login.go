package GUI

import (
	"ESFS2.0/client/common"
	"ESFS2.0/message/protos"
	"context"
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"time"
)

var _usernameEdit, _passwordEdit *walk.LineEdit
var _rememberRadioButton *walk.RadioButton

var HasLogin bool
var CurrentUser string
var Privatekeypath string

func init() {
	HasLogin = false //*********是否已登录********
	//CurrentUser = "memeshe"
}

func GetLoginPage() []Widget {
	var widget = []Widget{
		HSplitter{
			Children: []Widget{
				Label{
					Text: "用户名:",
				},
				LineEdit{
					AssignTo: &_usernameEdit,
				},
			},
		},
		HSplitter{
			Children: []Widget{
				Label{
					Text: "密码:",
				},
				LineEdit{
					AssignTo:     &_passwordEdit,
					PasswordMode: true,
				},
				RadioButton{
					AssignTo: &_rememberRadioButton,
					Text:     "记住用户名及密码",
				},
			},
		},
		PushButton{
			Text: "登录",
			OnClicked: func() {
				login()
			},
		},
		PushButton{
			Text: "注册",
			OnClicked: func() {
				CreateRegisterWindow()
			},
		},
	}
	return widget

}

func login() {
	c, conn, err := common.GetAuthenticationClient()
	if err != nil {
		log.Println(err)
		common.ShowMsgBox("提示", "服务器错误")
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	uname := _usernameEdit.Text()
	pwd := _passwordEdit.Text()

	if uname == "" || pwd == "" {
		common.ShowMsgBox("提示", "用户名和密码不能为空")
		return
	}

	request := &protos.LoginRequest{
		Username: uname,
		Password: pwd,
	}

	response, err := c.Login(ctx, request)

	switch response.ErrorMessage {
	case protos.ErrorMessage_SERVER_ERROR:
		common.ShowMsgBox("提示", "服务器错误")
	case protos.ErrorMessage_USER_NOT_EXISTS:
		common.ShowMsgBox("提示", "当前用户不存在")
	case protos.ErrorMessage_PASSWORD_WRONG:
		common.ShowMsgBox("提示", "密码错误")
	case protos.ErrorMessage_OK:
		common.ShowMsgBox("提示", "登录成功")
		HasLogin = true
		CurrentUser = uname
	}
}

func SelectPrivateKeyPath(mw *FileMainWindow) error {
	dlg := new(walk.FileDialog)
	dlg.Title = "请选择私钥的路径"
	dlg.Filter = "Certificate File(*.pem)|*.pem"
	if ok, err := dlg.ShowOpen(mw); !ok || err != nil {
		Privatekeypath = ""
		return errors.New("open path error")
	} else {
		Privatekeypath = dlg.FilePath
		return nil
	}
}
