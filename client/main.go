package main

import (
	"ESFS2.0/client/GUI"
	clicommon "ESFS2.0/client/common"
	"ESFS2.0/keyserver/common"
	"bytes"
	"path"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	walk.Resources.SetRootDirPath(path.Join(common.BaseDir, "client", "img"))

	mw := new(AppMainWindow)

	cfg := &GUI.MultiPageMainWindowConfig{
		Name:    "mainWindow",
		MinSize: Size{600, 400},
		MenuItems: []MenuItem{
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						Text:        "About",
						OnTriggered: func() { mw.aboutActionTriggered() },
					},
				},
			},
		},
		OnCurrentPageChanged: func() {
			mw.updateTitle(mw.CurrentPageTitle())
		},

		PageCfgs: []GUI.PageConfig{
			{"Login", "document-new.png", newLoginPage},
			{"Select", "document-properties.png", newSelectPage},
			//{"Upload", "document-properties.png", newUploadPage},
		},
	}

	mpmw, err := GUI.NewMultiPageMainWindow(cfg)
	if err != nil {
		panic(err)
	}

	mw.MultiPageMainWindow = mpmw

	mw.updateTitle(mw.CurrentPageTitle())

	mw.Run()
}

type AppMainWindow struct {
	*GUI.MultiPageMainWindow
}

func (mw *AppMainWindow) updateTitle(prefix string) {
	var buf bytes.Buffer

	if prefix != "" {
		buf.WriteString(prefix)
		buf.WriteString(" - ")
	}

	buf.WriteString("ESFS")

	_ = mw.SetTitle(buf.String())
}

func (mw *AppMainWindow) aboutActionTriggered() {
	walk.MsgBox(mw,
		"项目地址",
		"https://github.com/jiangxiaosheng/ESFS2.0.git",
		walk.MsgBoxOK|walk.MsgBoxIconInformation)
}

type LoginPage struct {
	*walk.Composite
}

func newLoginPage(parent walk.Container) (GUI.Page, error) {
	p := new(LoginPage)
	var widget []Widget
	if GUI.HasLogin == false {
		//未登录
		widget = GUI.GetLoginPage()
	} else {
		//已登录
		widget = []Widget{}
		var tmp walk.Form
		walk.MsgBox(tmp, "", "宁已登录", walk.MsgBoxIconInformation)
	}
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "loginPage",
		Layout:   VBox{},
		Children: widget,
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}

	return p, nil
}

type SelectPage struct {
	*walk.Composite
}

func newSelectPage(parent walk.Container) (GUI.Page, error) {
	p := new(SelectPage)
	var widget []Widget
	if GUI.HasLogin == true {
		//已登录
		widget = GUI.GetSelectPage()
	} else {
		//未登录
		widget = []Widget{}
		clicommon.ShowMsgBox("提示", "宁未登录")
	}
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "selectPage",
		Layout:   VBox{},
		Children: widget,
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}

	return p, nil
}
