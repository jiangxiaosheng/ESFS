package tests

import (
	"ESFS2.0/dataserver/common"
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	fmt.Println(common.GetConfigArray("db.tables")[0])
}
