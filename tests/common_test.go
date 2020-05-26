package tests

import (
	"ESFS2.0/utils"
	"fmt"
	"github.com/archivefile/zip"
	"io/ioutil"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	at := time.Now()
	fmt.Println(at.Format("2006-01-02 15:04:05"))
}

func TestXJB(t *testing.T) {
	b, err := utils.CompressToBytes("8.jpg")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("8.jpg.zip", b, 0644)
	//_, _ = utils.CompressToFile("8.jpg", "8.jpg.zip")
}

func TestZip(t *testing.T) {
	progress := func(archivePath string) {
		fmt.Println(archivePath)
	}
	err := zip.ArchiveFile("../message", "foo.zip", progress)
	if err != nil {
		fmt.Println(err)
	}
}
