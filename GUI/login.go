package main

import (
	"ESFS2.0/client"
	"ESFS2.0/message/protos"
	"context"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"time"
)

var _usernameEdit, _passwordEdit *walk.LineEdit
var _rememberRadioButton *walk.RadioButton

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
					AssignTo: &_passwordEdit,
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
				//fmt.Println(_usernameEdit.Text())
			},
		},
	}
	return widget
}

func login() {
	c, conn, err := client.GetAuthenticationClient()
	if err != nil {
		log.Println(err)
		ShowMsgBox("提示", "服务器错误")
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	uname := _usernameEdit.Text()
	pwd := _passwordEdit.Text()

	if uname == "" || pwd == "" {
		ShowMsgBox("提示", "用户名和密码不能为空")
		return
	}

	request := &protos.LoginRequest{
		Username: uname,
		Password: pwd,
	}

	response, err := c.Login(ctx, request)

	switch response.ErrorMessage {
	case protos.ErrorMessage_SERVER_ERROR:
		ShowMsgBox("提示", "服务器错误")
	case protos.ErrorMessage_USER_NOT_EXISTS:
		ShowMsgBox("提示", "当前用户不存在")
	case protos.ErrorMessage_PASSWORD_WRONG:
		ShowMsgBox("提示", "密码错误")
	case protos.ErrorMessage_OK:
		ShowMsgBox("提示", "登录成功")
		CurrentUser = uname
		HasLogin = true
	}
}
