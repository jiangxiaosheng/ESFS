package tests

import (
	"fmt"
	"github.com/archivefile/zip"
	"strings"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	at := time.Now()
	fmt.Println(at.Format("2006-01-02 15:04:05"))
}

func TestXJB(t *testing.T) {
	s := []string{"1", "2"}
	fmt.Println("(" + strings.Join(s, ",") + ")")
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
