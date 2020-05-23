package main

//
//import (
//	"github.com/lxn/walk"
//	. "github.com/lxn/walk/declarative"
//	"log"
//)
//
//var le1, le2 *walk.LineEdit
//var sport, maths, english *walk.RadioButton
//
//func main() {
//	if _, err := MainWindow1.Run(); err != nil {
//		log.Fatal(err)
//	}
//	log.Fatal(Bind("enabledCB.Checked"))
//}
//
//var MainWindow1 = MainWindow{
//	Title:    "Login",
//	MaxSize:  Size{100, 50},
//	Layout:   VBox{},
//	Children: widget,
//}
//
//var widget = []Widget{
//	HSplitter{
//		Children: []Widget{
//			lb1,
//			LineEdit1,
//		},
//		MaxSize: Size{5, 20},
//	},
//	HSplitter{
//		Children: []Widget{
//			lb2,
//			LineEdit2,
//			RadioButtonremember,
//		},
//		MaxSize: Size{5, 20},
//	},
//	HSplitter{
//		Children: []Widget{
//			PushButtonOK,
//		},
//		MaxSize: Size{5, 20},
//	},
//}
//var lb1 = Label{
//	Text: "用户名:",
//}
//var lb2 = Label{
//	Text: "密码:",
//}
//var LineEdit1 = LineEdit{
//	AssignTo: &le1,
//}
//var LineEdit2 = LineEdit{
//	AssignTo: &le2,
//}
//
//var RadioButtonremember = RadioButton{
//	AssignTo: &sport,
//	Text:     "记住用户名及密码",
//}
//var PushButtonOK = PushButton{
//	Text:      "登录",
//	OnClicked: OK_Clicked,
//}
//
//func OK_Clicked() {
//
//}
