// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"ESFS2.0/dataserver/common"
	"bytes"
	"path"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var HasLogin bool
var CurrentUser string

func init() {
	HasLogin = true //*********是否已登录********
	CurrentUser = "memeshe"
}

func main() {
	walk.Resources.SetRootDirPath(path.Join(common.BaseDir, "GUI", "img"))

	mw := new(AppMainWindow)

	cfg := &MultiPageMainWindowConfig{
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

		PageCfgs: []PageConfig{
			{"Login", "document-new.png", newLoginPage},
			{"Select", "document-properties.png", newSelectPage},
			//{"Upload", "document-properties.png", newUploadPage},
		},
	}

	mpmw, err := NewMultiPageMainWindow(cfg)
	if err != nil {
		panic(err)
	}

	mw.MultiPageMainWindow = mpmw

	mw.updateTitle(mw.CurrentPageTitle())

	mw.Run()
}

type AppMainWindow struct {
	*MultiPageMainWindow
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

func newLoginPage(parent walk.Container) (Page, error) {
	p := new(LoginPage)
	var widget []Widget
	if HasLogin == false {
		//未登录
		widget = GetLoginPage()
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

func newSelectPage(parent walk.Container) (Page, error) {
	p := new(SelectPage)
	var widget []Widget
	if HasLogin == true {
		//已登录
		widget = GetSelectPage()
	} else {
		//未登录
		widget = []Widget{}
		ShowMsgBox("提示", "宁未登录")
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
