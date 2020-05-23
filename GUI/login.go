package main

import (
	"fmt"
	"log"
	"os"
)
import (
	. "github.com/lxn/walk/declarative"
)

func init() {
	logFile, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Println(err.Error())
	}
	log.SetOutput(logFile)
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main() {
	if _, err := LoginWindow.Run(); err != nil {
		log.Fatal(err)
	}
	log.Fatal(Bind("enabledCB.Checked"))
}

var LoginWindow = MainWindow{
	Title:    "Login",
	MinSize:  Size{270, 290},
	Layout:   VBox{},
	Children: widget,
}

var widget = []Widget{
	HSplitter{
		Children: []Widget{
			usernameLabel,
			usernameEdit,
		},
		MaxSize: Size{50, 20},
	},
	HSplitter{
		Children: []Widget{
			passwordLabel,
			passwordEdit,
			RadioButtonRemember,
		},
		MaxSize: Size{50, 20},
	},
	PushButtonOK,
	PushButtonRegister,
}
var usernameLabel = Label{
	Text: "用户名:",
}
var passwordLabel = Label{
	Text: "密码:",
}

var usernameEdit = LineEdit{}

var passwordEdit = LineEdit{}

var RadioButtonRemember = RadioButton{
	Text: "记住用户名及密码",
}

var PushButtonOK = PushButton{
	Text:      "登录",
	OnClicked: okClicked,
}

var PushButtonRegister = PushButton{
	Text:      "注册",
	OnClicked: registerClicked,
}

func okClicked() {
	//client, conn, err := client.GetAuthenticationClient()
	//if err != nil {
	//	log.Printf(err.Error())
	//	return
	//}
	//defer conn.Close()
	//
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//
	//client.Login()

}
func registerClicked() {

}
