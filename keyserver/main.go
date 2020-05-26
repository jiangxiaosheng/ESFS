package keyserver

import (
	"ESFS2.0/message/protos"
)

type keyServer struct {
	protos.UnimplementedKeyStoreServer
}
