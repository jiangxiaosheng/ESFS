package tests

import (
	"ESFS2.0/utils"
	"testing"
)

func TestGenKeys(t *testing.T) {
	utils.GenerateRSAKey(1024)
}
