package main

import (
	"ESFS2.0/message"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
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
	m := new(FileModel)
	m.items = make([]*FileRecord, 2)

	m.items[0] = &FileRecord{
		FileInfo: message.FileInfo{
			Name:    "lalala",
			Mode:    0,
			Size:    0,
			ModTime: time.Time{},
		},
	}

	m.items[1] = &FileRecord{
		FileInfo: message.FileInfo{
			Name:    "hahaha",
			Mode:    0,
			Size:    0,
			ModTime: time.Time{},
		},
	}
	//******************展示示例：******************
	//m.items[0] = &FileRecord{
	//	Index:    0,
	//	FileInfo: message.FileInfo{
	//		Name:    "",
	//		Size:    0,
	//		ModTime: time.Time{},
	//	},
	//
	//}
	//

	//********************************************
	return m
}

type FileMainWindow struct {
	*walk.MainWindow
	model *FileModel
	tv    *walk.TableView
}

func GetSelectPage() []Widget {
	mw := &FileMainWindow{model: NewFileModel()}
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
						for _, x := range mw.model.items {
							if x.checked {
								//选中的文件
								fmt.Printf("checked: %v\n", x)
							}
						}
						fmt.Println()
					},
				},
				PushButton{
					Text: "上传",
					OnClicked: func() {
						OpenWindow()
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
						{Title: "日期"},
					},
					Model: mw.model,
					OnCurrentIndexChanged: func() {
						i := mw.tv.CurrentIndex()
						if 0 <= i {
							fmt.Printf("OnCurrentIndexChanged: %v\n", mw.model.items[i].Name)
						}
					},
					OnItemActivated: mw.tvItemActivated,
				},
			},
		},
	}
	return a
}

func logout() {
	HasLogin = false
	CurrentUser = ""
}

func listFiles() {

}

func (mw *FileMainWindow) tvItemActivated() {
	msg := ``
	for _, i := range mw.tv.SelectedIndexes() {
		msg = msg + "\n" + mw.model.items[i].Name
	}
	walk.MsgBox(mw, "title", msg, walk.MsgBoxIconInformation)
}
