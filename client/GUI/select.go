package GUI

import (
	clicommon "ESFS2.0/client/common"
	"ESFS2.0/message"
	"ESFS2.0/message/protos"
	"ESFS2.0/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io"
	"log"
	"net"
	"os"
	"path"
	"sort"
	"time"
)

type FileRecord struct {
	message.FileInfo
	checked bool
}

type FileModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*FileRecord
}

func (m *FileModel) RowCount() int {
	return len(m.items)
}

func (m *FileModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Name
	case 1:
		return item.Size
	case 2:
		return item.ModTime
	}
	panic("unexpected col")
}

//检查某一行是否选中
func (m *FileModel) Checked(row int) bool {
	return m.items[row].checked
}

//设置某一行的选中状态
func (m *FileModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked
	return nil
}

//根据某列对数据排序
func (m *FileModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.Stable(m)

	return m.SorterBase.Sort(col, order)
}

//获取数据行数
func (m *FileModel) Len() int {
	return len(m.items)
}

//内置的排序算法
func (m *FileModel) Less(i, j int) bool {
	a, b := m.items[i], m.items[j]

	c := func(ls bool) bool {
		if m.sortOrder == walk.SortAscending {
			return ls
		}

		return !ls
	}

	switch m.sortColumn {
	case 0:
		return c(a.Name < b.Name)
	case 1:
		return c(a.Size < b.Size)
	case 2:
		return c(a.ModTime.Before(b.ModTime))
	}

	panic("unreachable")
}

func (m *FileModel) Swap(i, j int) {
	m.items[i], m.items[j] = m.items[j], m.items[i]
}

func NewFileModel() *FileModel {
	c, conn, err := clicommon.GetFileHandleClient()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msg := &protos.ListFilesRequest{
		Username: CurrentUser,
	}

	response, err := c.ListFiles(ctx, msg)

	m := new(FileModel)
	m.items = []*FileRecord{}

	if response != nil {
		filesArray := response.FileInfo
		for _, data := range filesArray {
			info := &message.FileInfo{}
			err = json.Unmarshal(data, info)
			m.items = append(m.items, &FileRecord{
				FileInfo: *info,
				checked:  false,
			})
		}
	}

	return m
}

type FileMainWindow struct {
	*walk.MainWindow
	model *FileModel
	tv    *walk.TableView
}

func GetSelectPage() []Widget {
	mw := &FileMainWindow{}
	//异步渲染
	//go func() {
	//	mw.model = NewFileModel()
	//}()
	mw.model = NewFileModel()
	var a []Widget
	a = []Widget{
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				HSpacer{},
				Label{
					Text: "用户名:" + CurrentUser,
				},
				PushButton{
					Text: "下载",
					OnClicked: func() {
						selectDownloadFile(mw)
					},
				},
				PushButton{
					Text: "上传",
					OnClicked: func() {
						OpenWindow()
					},
				},
				PushButton{
					Text: "删除",
					OnClicked: func() {
						removeFiles(mw)
					},
				},
				PushButton{
					Text: "共享",
					OnClicked: func() {
						CreateShareWindow(mw)
					},
				},
				PushButton{
					Text: "登出",
					OnClicked: func() {
						logout()
					},
				},
			},
		},
		Composite{
			Layout: VBox{},
			ContextMenuItems: []MenuItem{
				Action{
					Text:        "I&nfo",
					OnTriggered: mw.tvItemActivated,
				},
				Action{
					Text: "E&xit",
					OnTriggered: func() {
						mw.Close()
					},
				},
			},
			Children: []Widget{
				TableView{
					AssignTo:         &mw.tv,
					CheckBoxes:       true,
					ColumnsOrderable: true,
					MultiSelection:   true,
					Columns: []TableViewColumn{
						{Title: "名称"},
						{Title: "大小"},
						{Title: "日期", FormatFunc: func(t interface{}) string {
							return fmt.Sprintf(t.(time.Time).Format("2006-01-02 15:04:05"))
						}},
					},
					Model:           mw.model,
					OnItemActivated: mw.tvItemActivated,
				},
			},
		},
	}
	MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "文件上传",
		MinSize:  Size{300, 400},
		Layout:   VBox{},
		Children: []Widget{},
		Visible:  false,
	}.Create()
	mw.MainWindow.Close()

	return a
}

func selectDownloadFile(mw *FileMainWindow) {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	if ok, err := dlg.ShowBrowseFolder(mw); err != nil {
		return
	} else if !ok {
		return
	} else {
		fmt.Println(dlg.FilePath)
		go func() {
			download(mw, dlg.FilePath)
		}()
	}
}

/**
@author js
下载云上文件
*/
func download(mw *FileMainWindow, dir string) {
	fileItems := mw.model.items
	var filesToDownload []string
	for _, fileRecord := range fileItems {
		if fileRecord.checked {
			filesToDownload = append(filesToDownload, fileRecord.Name)
		}
	}

	//1.使用grpc获取文件-二级密码map
	c, conn, err := clicommon.GetFileHandleClient()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &protos.DownloadPrepareRequest{
		Username: CurrentUser,
	}
	response, err := c.DownloadPrepare(ctx, request)
	if err != nil {
		log.Printf("获取二级密码失败 %v", err.Error())
		clicommon.ShowMsgBox("提示", "服务器错误")
		return
	}
	var m map[string]string

	err = json.Unmarshal(response.Content, &m)
	if err != nil {
		log.Printf("反序列化失败 %v", err.Error())
		clicommon.ShowMsgBox("提示", "服务器错误")
		return
	}

	//2.获取当前用户的私钥，用来后续解密文件
	priKey := clicommon.GetUserPrivateKey()
	if priKey == nil {
		log.Printf("读取私钥失败 %v", err.Error())
		clicommon.ShowMsgBox("提示", "服务器错误")
		return
	}

	//3.建立socket连接，开始传输数据，这部分的实现比较tricky，可以不用深入研究
	msg := message.FileSocketMessage{
		UserName: CurrentUser,
		FileName: filesToDownload,
		Type:     message.FILE_DOWNLOAD,
	}
	addr := "0.0.0.0:8959"
	socketConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("建立socket连接失败 %v", err)
		clicommon.ShowMsgBox("提示", "服务器错误")
		return
	}
	serializedData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("序列化失败 %v", err.Error())
		clicommon.ShowMsgBox("提示", "服务器错误")
		return
	}
	_, err = socketConn.Write(serializedData)
	if err != nil {
		log.Printf("socket写入数据失败 %v", err.Error())
		clicommon.ShowMsgBox("提示", "服务器错误")
		return
	}

	buffer := make([]byte, 2048)
	socketConn.Read(buffer)
	socketConn.Write([]byte{1})

	for _, filename := range filesToDownload {
		signal := &message.SignalOver{}
		var encryptedData []byte
		for {
			n, err := socketConn.Read(buffer)
			if err == io.EOF {
				break
			}
			if json.Unmarshal(buffer[:n], signal) == nil {
				break
			}

			encryptedData = append(encryptedData, buffer[:n]...)
			if err != nil {
				log.Printf("读取字节失败 %v", err.Error())
				clicommon.ShowMsgBox("提示", "服务器错误")
				return
			}
			socketConn.Write([]byte{1})
		}
		//获得加密文件时用的会话密钥
		key, err := utils.GenerateSessionKeyWithSecondKey(m[filename], priKey)
		if err != nil {
			log.Printf("生成会话密钥失败 %v", err.Error())
			clicommon.ShowMsgBox("提示", "服务器错误")
			return
		}
		//使用会话密钥解密云上文件
		err = utils.AESDecryptToFile(encryptedData, key, path.Join(dir, filename))
		if err != nil {
			clicommon.ShowMsgBox("提示", "服务器错误")
			return
		}

		socketConn.Write([]byte{1})

		//下载签名文件，用来进行认证
		sigFile, err := os.Create(path.Join(dir, "."+filename+".sig"))
		if err != nil {
			log.Printf("签名文件创建失败 %v", err.Error())
			clicommon.ShowMsgBox("提示", "服务器错误")
			return
		}

		log.Println("创建签名文件")
		for {
			n, err := socketConn.Read(buffer)
			if err == io.EOF {
				break
			}
			fmt.Println(string(buffer[:n]))
			if json.Unmarshal(buffer[:n], signal) == nil {
				break
			}
			_, err = sigFile.Write(buffer[:n])
			if err != nil {
				log.Printf("签名文件写入失败 %v", err.Error())
				clicommon.ShowMsgBox("提示", "服务器错误")
				return
			}
			socketConn.Write([]byte{1})
		}
		log.Println("签名文件写入完毕")
		sigFile.Close()
		socketConn.Write([]byte{1})
	}
	clicommon.ShowMsgBox("提示", "下载成功")
}

func logout() {
	HasLogin = false
	CurrentUser = ""
}

func (mw *FileMainWindow) tvItemActivated() {
	msg := ``
	for _, i := range mw.tv.SelectedIndexes() {
		msg = msg + "\n" + mw.model.items[i].Name
	}
	walk.MsgBox(mw, "title", msg, walk.MsgBoxIconInformation)
}

/**
@author js
删除服务器上的指定文件
*/
func removeFiles(mw *FileMainWindow) {
	fileItems := mw.model.items
	var filesToRemove []string
	for _, fileRecord := range fileItems {
		if fileRecord.checked {
			filesToRemove = append(filesToRemove, fileRecord.Name)
		}
	}

	go func() {
		c, conn, err := clicommon.GetFileHandleClient()
		if err != nil {
			fmt.Println(err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		request := &protos.RemoveFilesRequest{
			Username:  CurrentUser,
			Filenames: filesToRemove,
		}

		response, err := c.RemoveFiles(ctx, request)
		if response.ErrorMessage != protos.ErrorMessage_OK {
			clicommon.ShowMsgBox("提示", "删除失败")
			return
		}
		mw.model = NewFileModel()
		clicommon.ShowMsgBox("提示", "删除成功")
	}()

}
