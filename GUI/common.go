package main

import "github.com/lxn/walk"

func ShowMsgBox(title, content string) {
	var tmp walk.Form
	walk.MsgBox(tmp, title, content, walk.MsgBoxIconInformation)
}
