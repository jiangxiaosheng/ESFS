package GUI

import (
	clicommon "ESFS2.0/client/common"
	"ESFS2.0/message"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type MyMainWindow struct {
	*walk.MainWindow
	selectedfile *walk.LineEdit
	secondpw     *walk.LineEdit
	message      *walk.ListBox
	wv           *walk.WebView
	textEdit     *walk.TextEdit
}

var _useRadioButton, _notuseRadioButton *walk.RadioButton

func OpenWindow() {
	mw := &MyMainWindow{}
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "文件上传",

		MinSize: Size{300, 400},
		Layout:  VBox{},
		Children: []Widget{
			GroupBox{
				MaxSize: Size{500, 500},
				Layout:  HBox{},
				Children: []Widget{
					PushButton{ //选择文件
						Text:      "打开文件",
						OnClicked: mw.selectFile, //点击事件
					},
					Label{Text: "选中的文件 "},
					LineEdit{
						AssignTo: &mw.selectedfile, //选中的文件
					},
					PushButton{
						Text: "上传",
						OnClicked: func() {
							upload(mw)
						},
					},
				},
			},
			GroupBox{
				MaxSize: Size{500, 500},
				Layout:  HBox{},
				Children: []Widget{
					RadioButton{
						AssignTo: &_useRadioButton,
						Text:     "使用默认二级密码",
						OnClicked: func() {
							mw.secondpw.SetReadOnly(true)
						},
					},
					RadioButton{
						AssignTo: &_notuseRadioButton,
						Text:     "使用新二级密码",
						OnClicked: func() {
							mw.secondpw.SetReadOnly(false)
						},
					},
					LineEdit{
						AssignTo: &mw.secondpw, //新的二级密码
						ReadOnly: true,
					},
				},
			},
			TextEdit{
				//MaxSize: Size{300, 200},
				AssignTo: &mw.textEdit,
				ReadOnly: true,
				Text:     "Drop files here, from windows explorer...",
			},
		},
	}.Create()); err != nil {
		fmt.Printf("Run err: %+v\n", err)
	}
	mw.textEdit.DropFiles().Attach(func(files []string) {
		mw.textEdit.SetText(strings.Join(files, "\r\n"))
	})
	mw.Run()
}

func upload(mw *MyMainWindow) {
	path := mw.selectedfile.Text()
	file, err := os.Open(path)
	if err != nil {
		log.Printf("打开文件失败 %v", err.Error())
		clicommon.ShowMsgBox("提示", "打开文件失败")
		return
	}

	addr := "0.0.0.0:8959"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	msg := message.FileSocketMessage{
		UserName: CurrentUser,
		FileName: filepath.Base(file.Name()),
		Type:     message.FILE_UPLOAD,
	}

	serializedData, err := json.Marshal(msg)
	_, err = conn.Write(serializedData)
	if err != nil {
		fmt.Printf("socket写入数据失败 %v", err.Error())
		return
	}

	buffer := make([]byte, 2048)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	clicommon.ShowMsgBox("提示", "上传成功")

	//异步上传
	go func() {

	}()
}

func (mw *MyMainWindow) selectFile() {
	allowType := viper.GetStringSlice("common.server.allowtype")
	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	filter := []string{}
	filterstring := ""
	for _, v := range allowType {
		if v != "*" {
			filterstring = "可上传" + v + " (*." + v + ")|*." + v
			filter = append(filter, filterstring)
		} else {
			filterstring = "所有文件" + v + " (*." + v + ")|*." + v
			filter = append(filter, filterstring)
		}
	}
	dlg.Filter = strings.Join(filter, "|") //切片转换字符串
	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.selectedfile.SetText("") //通过重定向变量设置TextEdit的Text
		return
	} else if !ok {
		mw.selectedfile.SetText("") //通过重定向变量设置TextEdit的Text
		return
	}
	mw.selectedfile.SetText(dlg.FilePath) //通过重定向变量设置TextEdit的Text
}
