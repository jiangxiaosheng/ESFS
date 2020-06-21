package common

import (
	"ESFS2.0/keyserver/common"
	"ESFS2.0/message/protos"
	"ESFS2.0/utils"
	"crypto/rsa"
	"fmt"
	"github.com/lxn/walk"
	"google.golang.org/grpc"
	"log"
	"path"
)

func ShowMsgBox(title, content string) {
	var tmp walk.Form
	walk.MsgBox(tmp, title, content, walk.MsgBoxIconInformation)
}

func GetAuthenticationClient() (protos.AuthenticationClient, *grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8927)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("无法建立grpc连接 %v", err.Error())
		return nil, conn, err
	}

	c := protos.NewAuthenticationClient(conn)
	return c, conn, nil
}

func GetFileHandleClient() (protos.FileHandleClient, *grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8927)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("无法建立grpc连接 %v", err.Error())
		return nil, conn, err
	}

	c := protos.NewFileHandleClient(conn)
	return c, conn, nil
}

func GetCAClient() (protos.CAClient, *grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 9015)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("无法建立grpc连接 %v", err.Error())
		return nil, conn, err
	}

	c := protos.NewCAClient(conn)
	return c, conn, nil
}

func GetUserPrivateKey(Privatekeypath string) *rsa.PrivateKey {
	keyPath := Privatekeypath
	log.Println(keyPath)
	//path.Join(common.BaseDir, "client", "keys", "private.pem")/*hard code version*/
	key := utils.GetPrivateKeyFromFile(keyPath)
	return key
}

func GetUserPublicKey() *rsa.PublicKey {
	keyPath := path.Join(common.BaseDir, "client", "keys", "public.pem")
	key := utils.GetPublicKeyFromFile(keyPath)
	return key
}
