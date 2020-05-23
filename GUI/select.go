package main

//
//import (
//	"fmt"
//	"github.com/lxn/walk"
//	. "github.com/lxn/walk/declarative"
//	"sort"
//)
//
//type File struct {
//	Index    int
//	Name     string
//	Size     int
//	Filetime int
//	checked  bool
//}
//
//type FileModel struct {
//	walk.TableModelBase
//	walk.SorterBase
//	sortColumn int
//	sortOrder  walk.SortOrder
//	items      []*File
//}
//
//func (m *FileModel) RowCount() int {
//	return len(m.items)
//}
//
//func (m *FileModel) Value(row, col int) interface{} {
//	item := m.items[row]
//
//	switch col {
//	case 0:
//		return item.Index
//	case 1:
//		return item.Name
//	case 2:
//		return item.Size
//	case 3:
//		return item.Filetime
//	}
//	panic("unexpected col")
//}
//
//func (m *FileModel) Checked(row int) bool {
//	return m.items[row].checked
//}
//
//func (m *FileModel) SetChecked(row int, checked bool) error {
//	m.items[row].checked = checked
//	return nil
//}
//
//func (m *FileModel) Sort(col int, order walk.SortOrder) error {
//	m.sortColumn, m.sortOrder = col, order
//
//	sort.Stable(m)
//
//	return m.SorterBase.Sort(col, order)
//}
//
//func (m *FileModel) Len() int {
//	return len(m.items)
//}
//
//func (m *FileModel) Less(i, j int) bool {
//	a, b := m.items[i], m.items[j]
//
//	c := func(ls bool) bool {
//		if m.sortOrder == walk.SortAscending {
//			return ls
//		}
//
//		return !ls
//	}
//
//	switch m.sortColumn {
//	case 0:
//		return c(a.Index < b.Index)
//	case 1:
//		return c(a.Name < b.Name)
//	case 2:
//		return c(a.Size < b.Size)
//	case 3:
//		return c(a.Filetime < b.Filetime)
//	}
//
//	panic("unreachable")
//}
//
//func (m *FileModel) Swap(i, j int) {
//	m.items[i], m.items[j] = m.items[j], m.items[i]
//}
//
//func NewFileModel() *FileModel {
//	m := new(FileModel)
//	m.items = make([]*File, 3)
//	//******************展示示例：******************
//	m.items[0] = &File{
//		Index:    0,
//		Name:     "杜蕾斯",
//		Size:     20,
//		Filetime: 20171020122634, //2017.10.20.12.26.34
//	}
//
//	m.items[1] = &File{
//		Index:    1,
//		Name:     "杰士邦",
//		Size:     18,
//		Filetime: 20171020122634,
//	}
//
//	m.items[2] = &File{
//		Index:    2,
//		Name:     "冈本",
//		Size:     19,
//		Filetime: 20171020122634,
//	}
//	//********************************************
//	return m
//}
//
//type FileMainWindow struct {
//	*walk.MainWindow
//	model *FileModel
//	tv    *walk.TableView
//}
//
//func main() {
//	mw := &FileMainWindow{model: NewFileModel()}
//
//	MainWindow{
//		AssignTo: &mw.MainWindow,
//		Title:    "文件展示",
//		Size:     Size{800, 600},
//		Layout:   VBox{},
//		Children: []Widget{
//			Composite{
//				Layout: HBox{MarginsZero: true},
//				Children: []Widget{
//					HSpacer{},
//					Label{
//						Text: "用户名:ytw",
//					},
//					Label{
//						Text: "邮箱:123456@qq.com",
//					},
//					PushButton{
//						Text: "下载",
//						OnClicked: func() {
//							for _, x := range mw.model.items {
//								if x.checked {
//									//选中的文件
//									fmt.Printf("checked: %v\n", x)
//								}
//							}
//							fmt.Println()
//						},
//					},
//					PushButton{
//						Text: "上传",
//						OnClicked: func() {
//
//						},
//					},
//					PushButton{
//						Text: "登出",
//						OnClicked: func() {
//
//						},
//					},
//				},
//			},
//			Composite{
//				Layout: VBox{},
//				ContextMenuItems: []MenuItem{
//					Action{
//						Text:        "I&nfo",
//						OnTriggered: mw.tv_ItemActivated,
//					},
//					Action{
//						Text: "E&xit",
//						OnTriggered: func() {
//							mw.Close()
//						},
//					},
//				},
//				Children: []Widget{
//					TableView{
//						AssignTo:         &mw.tv,
//						CheckBoxes:       true,
//						ColumnsOrderable: true,
//						MultiSelection:   true,
//						Columns: []TableViewColumn{
//							{Title: "编号"},
//							{Title: "名称"},
//							{Title: "大小"},
//							{Title: "日期"},
//						},
//						Model: mw.model,
//						OnCurrentIndexChanged: func() {
//							i := mw.tv.CurrentIndex()
//							if 0 <= i {
//								fmt.Printf("OnCurrentIndexChanged: %v\n", mw.model.items[i].Name)
//							}
//						},
//						OnItemActivated: mw.tv_ItemActivated,
//					},
//				},
//			},
//		},
//	}.Run()
//}
//
//func (mw *FileMainWindow) tv_ItemActivated() {
//	msg := ``
//	for _, i := range mw.tv.SelectedIndexes() {
//		msg = msg + "\n" + mw.model.items[i].Name
//	}
//	walk.MsgBox(mw, "title", msg, walk.MsgBoxIconInformation)
//}
