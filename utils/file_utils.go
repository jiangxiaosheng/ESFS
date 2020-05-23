package utils

import (
	"archive/zip"
	"crypto/rsa"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	basedir = "E:\\GoLand\\GoLand 2019.3.3\\codes\\src\\ESFS"
)

func GetDir(dir string) string {
	return path.Join(basedir, dir)
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
				var rbase = base + file.Name() + "/"
				Rfile, error := os.Open(base + file.Name())
				if error != nil {
					println(error)
				}
				*fileArr = append(*fileArr, Rfile)
				rev(rbase, fileArr)
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
			Rfile, error := os.Open(dir)
			if error != nil {
				println(error)
			} else {
				*fileArr = append(*fileArr, Rfile)
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
TODO:不能用，想写一个
*/
func UploadFiles(path string, priKey *rsa.PrivateKey, secondKey string) {
	files := new([]*os.File)
	ReadFiles(path, files)
	//comFile, _ := Compress(*files, "test.zip")
	//fmt.Println(comFile.Name())
	for _, i := range *files {
		fmt.Println(i.Name())
	}
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

//////////////////下面代码都是网上扒的，还没改好

//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Compress(files []*os.File, dest string) (*os.File, error) {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//解压
func Decompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

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

////////////////////////////////////////////////////////
