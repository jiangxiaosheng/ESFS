package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	at := time.Now()
	fmt.Println(at.Format("2006-01-02 15:04:05"))
}

func TestXJB(t *testing.T) {
	file, _ := os.Open("C:\\Users\\jiangsheng\\Desktop\\rad.jpg")
	fmt.Println(filepath.Base(file.Name()))
}
