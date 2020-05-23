package GUI

import (
	"fmt"
	"log"
)
import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var _usernameEdit, _passwordEdit *walk.LineEdit
var _rememberRadioButton *walk.RadioButton

var loginWindow = MainWindow{
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
			rememberPushButton,
		},
		MaxSize: Size{50, 20},
	},
	okPushButton,
	registerPushButton,
}

var usernameLabel = Label{
	Text: "用户名:",
}

var passwordLabel = Label{
	Text: "密码:",
}

var usernameEdit = LineEdit{
	AssignTo: &_usernameEdit,
}

var passwordEdit = LineEdit{
	AssignTo: &_passwordEdit,
}

var rememberPushButton = RadioButton{
	AssignTo: &_rememberRadioButton,
	Text:     "记住用户名及密码",
}

var okPushButton = PushButton{
	Text:      "登录",
	OnClicked: okClicked,
}
var registerPushButton = PushButton{
	Text:      "注册",
	OnClicked: registerClicked,
}

func okClicked() {
	fmt.Println(_usernameEdit.Text())
}

func registerClicked() {

}

func main() {
	if _, err := loginWindow.Run(); err != nil {
		log.Fatal(err)
	}
	log.Fatal(Bind("enabledCB.Checked"))
}
