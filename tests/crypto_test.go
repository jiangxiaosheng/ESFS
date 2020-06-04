package tests

import (
	"ESFS2.0/client/common"
	"ESFS2.0/utils"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGenKeys(t *testing.T) {
	utils.GenerateRSAKey(1024, "./", "")
}

func TestDS(t *testing.T) {
	priKey := common.GetUserPrivateKey()
	file := "C:\\Users\\jiangsheng\\Desktop\\rad.jpg"
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print(err)
	}
	b, err := utils.GenerateDS(data, priKey)
	if err != nil {
		fmt.Print(err)
	}
	ioutil.WriteFile("testsig.txt", b, 0644)
	fmt.Println(b)
}

func TestAES(t *testing.T) {
	var b []byte
	a := []byte{1, 2, 3}
	b = append(b, a...)
	println(b)
}
