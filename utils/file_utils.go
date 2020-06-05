package utils

import (
	"bytes"
	"fmt"
	"github.com/archivefile/zip"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/**
@author yyx
DONE: path为需要压缩的路径地址（不需要判断是否为目录）,dest为输出路径，返回相应的压缩文件和错误，d为返回的*os.File对象
*/
func CompressToFile(path string, dest string) (*os.File, error) {
	err := zip.ArchiveFile(path, dest, nil)
	if err != nil {
		log.Printf("压缩文件失败 %v", err.Error())
		return nil, err
	}

	d, err := os.Open(dest)
	if err != nil {
		log.Printf("打开文件出错")
		return nil, err
	}
	return d, nil
}

/**
@author yyx
*/
func CompressToBytes(path string) ([]byte, error) {
	var b []byte
	writer := bytes.NewBuffer(b)

	err := zip.Archive(path, writer, nil)
	if err != nil {
		log.Printf("压缩文件失败 %v", err.Error())
		return nil, err
	}

	writer.Bytes()
	return writer.Bytes(), nil
}

/**
@author: yyx
*/
func rev(base string, fileArr *[]*os.File) {
	files, err := ioutil.ReadDir(base)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				Rfile, error := os.Open(base + file.Name())
				if error != nil {
					println("Wrong")
				}
				*fileArr = append(*fileArr, Rfile)
			} else {
				//var rbase = base + file.Name() + "/"
				Rfile, error := os.Open(base + file.Name())
				if error != nil {
					println(error)
				}
				*fileArr = append(*fileArr, Rfile)
				//rev(rbase, fileArr)
			}
		}
	} else {
		println("Opendir Failed!")
	}
}

/**
@author: yyx
*/
func ReadFiles(dir string, fileArr *[]*os.File) {
	file, err := os.Stat(dir)
	if err == nil {
		if !file.IsDir() {
			Rfile, error := os.Open(file.Name())
			if error != nil {
				println("Wrong")
			}
			*fileArr = append(*fileArr, Rfile)
		} else {
			var base = dir + "/"
			_, error := os.Open(dir)
			if error != nil {
				println(error)
			} else {
				//*fileArr = append(*fileArr, Rfile)
				rev(base, fileArr)
			}
		}
	} else {
		fmt.Println(err.Error())
		panic(err)
	}
	return
}

/**
文件转换为字节数组
*/
func FileToBytes(file *os.File) []byte {
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	return buf
}

/**
@author ytw
读取文件
*/
func ReadFile(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

/**
@author ytw
写入文件
*/
func WriteFile(filePth string, data []byte) error {
	err := ioutil.WriteFile(filePth, data, 0644)
	return err
}

//func EncryptFile(file *os.File, key )

//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
func main() { //测试
	var list []*os.File
	ReadFiles("C:/Users/Administrator/go/copy", &list)
	for i := 0; i < len(list); i++ {
		//fmt.Println(list[i].Name())
	}
	//CompressToFile("C:/Users/Administrator/go/src", "C:/Users/Administrator/go/copy/new.zip")
	CompressToBytes("C:/Users/Administrator/go/src")
	//str := string(b)
}

////////////////////////////////////////////////////////
