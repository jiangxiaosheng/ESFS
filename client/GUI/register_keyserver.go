package GUI

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MyKeyserverRegisterWindow struct {
	*walk.MainWindow
	username    *walk.LineEdit
	selectcer   *walk.LineEdit
	generatecer *walk.LineEdit
}

func kregist() {}
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
					Label{
						Text: "选择路径",
					},
					LineEdit{
						AssignTo: &mw.selectcer,
					},
					PushButton{
						Text: "选择证书",
						OnClicked: func() {
							mw.SelectCer() //TODO:选择证书
						},
					},
				},
			},
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "选择路径",
					},
					LineEdit{
						AssignTo: &mw.generatecer,
					},
					PushButton{
						Text: "生成证书",
						OnClicked: func() {
							mw.GenerateCer() //TODO:生成证书
						},
					},
				},
			},
			PushButton{
				Text: "注册",
				OnClicked: func() {
					//kregist()//TODO:keyserver注册
				},
			},
		},
	}.Create()); err != nil {
		fmt.Println("Error!")
	}
}

func (mw *MyKeyserverRegisterWindow) SelectCer() {
	Dlg := new(walk.FileDialog)
	Dlg.Title = "选择证书"
	Dlg.Filter = "Certificate File(*.cer)|*.cer" //
	if ok, err := Dlg.ShowOpen(mw); err != nil || !ok {
		mw.selectcer.SetText("")
		return
	} else {
		mw.selectcer.SetText(Dlg.FilePath)
		return
	}
}

func (mw *MyKeyserverRegisterWindow) GenerateCer() {
	Dlg := new(walk.FileDialog)
	Dlg.Title = "选择生成路径"
	Dlg.Filter = "Certificate File(*.cer)|*.cer" //
	Dlg.InitialDirPath = "Default"
	if ok, err := Dlg.ShowSave(mw); err != nil || !ok {
		mw.selectcer.SetText("")
		return
	} else {
		mw.selectcer.SetText(Dlg.FilePath)
	}
	return
}
