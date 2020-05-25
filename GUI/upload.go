package main

import (
	"ESFS2.0/message"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type MyMainWindow struct {
	*walk.MainWindow
	model        *MessageModel
	selectedfile *walk.LineEdit
	message      *walk.ListBox
	wv           *walk.WebView
	textEdit     *walk.TextEdit
}

func OpenWindow() {
	mw := &MyMainWindow{model: NewMessageModel()}

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
			ListBox{ //记录框
				AssignTo:              &mw.message,
				OnCurrentIndexChanged: mw.lb_CurrentIndexChanged, //单击
				OnItemActivated:       mw.lb_ItemActivated,       //双击
			},
			TextEdit{
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
		mw.selectedfile.SetText(strings.Join(files, "\r\n"))
	})
	mw.Run()
}

func upload(mw *MyMainWindow) {
	path := mw.selectedfile.Text()
	file, err := os.Open(path)
	if err != nil {
		log.Printf("打开文件失败 %v", err.Error())
		ShowMsgBox("提示", "打开文件失败")
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

	ShowMsgBox("提示", "上传成功")

	//异步上传
	go func() {

	}()
}

func (mw *MyMainWindow) selectFile() {
	allowType := viper.GetStringSlice("common.server.allowtype")
	fmt.Printf("allowType: %+v\n", allowType)

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	//dlg.Filter = "可上传jpg (*.jpg)|*.jpg|可上传png (*.png)|*.png|可上传gif (*.gif)|*.gif|所有文件 (*.*)|*.*"
	//判断可允许上传的文件
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
	fmt.Printf("dlg.Filter: %+v\n", dlg.Filter)

	record := getMessage(mw) //获取记录
	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.selectedfile.SetText("")                       //通过重定向变量设置TextEdit的Text
		_ = writeMessage(mw, "Error : File Open", record) //写入记录
		return
	} else if !ok {
		mw.selectedfile.SetText("")            //通过重定向变量设置TextEdit的Text
		_ = writeMessage(mw, "cancel", record) //写入记录
		return
	}
	s := fmt.Sprintf("Select : %s", dlg.FilePath)
	_ = writeMessage(mw, s, record)       //写入记录
	mw.selectedfile.SetText(dlg.FilePath) //通过重定向变量设置TextEdit的Text

	//************************文件path**************************************************
	fmt.Printf("dlg.FilePath: %+v\n", dlg.FilePath)
	//********************************************************************************
}

func writeMessage(mw *MyMainWindow, message string, record []string) []string {
	is_http := strings.Index(message, "http") //查找字符串位置
	message_record := message
	if is_http != -1 { // -1 是找不到
		message_record = "双击查看图片：" + message

		//插入默认浏览器打开记录
		item := MessageItem{
			Name:  "双击用默认浏览器打开图片" + message,
			Value: message,
		}
		mw.model.items = append(mw.model.items, item)
		record = append(record, "双击用默认浏览器打开图片"+message) //插入记录
	}
	record = append(record, message_record) //插入记录
	appendMessageModel(mw, message)         //插入模型
	mw.message.SetModel(record)             //记录输出

	return record
}

//消息记录改变事件
func (mw *MyMainWindow) lb_CurrentIndexChanged() {
	fmt.Printf("mw.message.CurrentIndex(): ", mw.message.CurrentIndex())
	fmt.Println()
	return
}

//消息记录点击事件
func (mw *MyMainWindow) lb_ItemActivated() {
	fmt.Println("mw.message.CurrentIndex(): ", mw.message.CurrentIndex())
	fmt.Println()
	fmt.Printf("mw.model.items: %+v ", mw.model.items)
	fmt.Println()

	index := mw.message.CurrentIndex()
	imagename := mw.model.items[index].Name    //获取当前选中名称
	image := mw.model.items[index].Value       //获取当前选中值
	is_http := strings.Index(image, "http")    //查找字符串位置
	is_ie := strings.Index(imagename, "默认浏览器") //查找字符串位置
	if is_http != -1 {                         // -1 是找不到
		//walk.MsgBox(mw, "Value", value, walk.MsgBoxIconInformation) //提示框
		fmt.Printf("image : %+v ", image)
		fmt.Println()

		if is_ie != -1 {
			openImageExplorer(image) ////使用默认浏览器打开图片
		} else {
			openImageWebview(mw, image) ////使用webview打开图片
		}
	}
	return
}

//使用默认浏览器打开图片
func openImageExplorer(image string) {
	cmd := exec.Command("explorer", image)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	return
}

//创建html,使用webview打开图片
func openImageWebview(mw *MyMainWindow, image string) {
	go func() {
		h := md5.New()
		io.WriteString(h, image)
		htmlname := fmt.Sprintf("%x", h.Sum(nil)) //生成html名称
		userFile := "" + htmlname + ".html"
		fileall := getCurrentDirectory() + "/" + userFile //html的绝对路径

		fout, err := os.Create(fileall) //创建html文件
		defer fout.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		html := `<!DOCTYPE html><html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8" /><title>图片展示</title></head><body><a href="%s" target="_self">%s</a><br><img src="%s" alt=""></body></html>`
		html = fmt.Sprintf(html, image, image, image) //插入图片
		fout.WriteString(html)                        //把字符串写入html

		mw.wv.SetURL("file:///" + getCurrentDirectory() + "/" + userFile)
	}()
}

//获取当前文件路径
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//消息记录每个模型
type MessageItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//消息记录模型
type MessageModel struct {
	walk.ListModelBase
	items []MessageItem
}

//新建消息模型
func NewMessageModel() *MessageModel {
	m := &MessageModel{items: []MessageItem{}}
	return m
}

//插入到消息模型中
func appendMessageModel(mw *MyMainWindow, message string) {
	item := MessageItem{
		Name:  message,
		Value: message,
	}
	mw.model.items = append(mw.model.items, item)

}

//获取记录
func getMessage(mw *MyMainWindow) []string {
	message := mw.message.Model() //获取以前的记录
	fmt.Println("message", message)
	record := []string{} //记录
	if message != nil {
		for _, v := range message.([]string) {
			record = append(record, v) //插入以前的记录
		}
	}
	return record
}
