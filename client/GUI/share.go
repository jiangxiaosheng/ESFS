package GUI

import (
	clicommon "ESFS2.0/client/common"
	"ESFS2.0/message/protos"
	"ESFS2.0/utils"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"time"
)

type MyShareWindow struct {
	*walk.MainWindow
	shareName *walk.LineEdit
}

func CreateShareWindow(father *FileMainWindow) {
	mw := &MyShareWindow{}
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "文件分享",
		Size:     Size{450, 100},
		Layout:   VBox{},
		Children: []Widget{
			GroupBox{
				//MaxSize: Size{500, 500},
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "用户名"},
					LineEdit{
						AssignTo: &mw.shareName,
					},
				},
			},
			PushButton{
				Text: "确定",
				OnClicked: func() {
					share(father, mw)
				},
			},
		},
	}.Create()); err != nil {
		fmt.Println("Error!")
	}
}

func share(father *FileMainWindow, mw *MyShareWindow) {
	if mw.shareName.Text() == "" {
		clicommon.ShowMsgBox("提示", "请输入用户名")
		return
	}

	exists, err := checkUserExists(mw.shareName.Text())
	if err != nil {
		clicommon.ShowMsgBox("提示", "服务器错误")
		return
	}
	if exists == false {
		clicommon.ShowMsgBox("提示", "用户不存在")
		return
	}
	//priKey := clicommon.GetUserPrivateKey()

	//rsa.EncryptPKCS1v15(rand.Reader, key, data)

}

/**
@author yyx
返回文件名-二级密码映射
*/
func getSecondKeys(username string, filenames []string) (map[string]string, error) {
	c, conn, err := clicommon.GetAuthenticationClient()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := &protos.GetSecondKeyRequest{
		Username:  username,
		Filenames: filenames,
	}

	response, err := c.GetSecondKey(ctx, request)
	switch response.ErrorMessage {
	case protos.ErrorMessage_SERVER_ERROR:
		return nil, errors.New("server_error")
	case protos.ErrorMessage_USER_NOT_EXISTS:
		return nil, errors.New("username_not_exist")
	case protos.ErrorMessage_OK:
		var fsmap = make(map[string]string)
		err = json.Unmarshal(response.SecondKeysMapData, fsmap)
		if err != nil {
			return nil, err
		}
		return fsmap, nil
	}
	return nil, nil
}

/**
@author js
检查某个用户是否存在，不存在就不能共享了
*/
func checkUserExists(username string) (bool, error) {
	c, conn, err := clicommon.GetAuthenticationClient()
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &protos.CheckUserExistsRequest{Username: username}
	response, err := c.CheckUserExits(ctx, request)
	if response.ErrorMessage == protos.ErrorMessage_SERVER_ERROR {
		return true, errors.New("服务器错误")
	}
	return response.Exists, nil
}

/**
@author js
获取想要共享的那个用户的公钥（该公钥用来对二级密码加密后添加到access数据表中）
*/
func getUserPublicKey(username string) (*rsa.PublicKey, protos.ErrorMessage) {
	c, conn, err := clicommon.GetCAClient()
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &protos.GetCertRequest{
		Username: username,
	}
	response, err := c.GetCert(ctx, request)
	switch response.ErrorMessage {
	case protos.ErrorMessage_SERVER_ERROR:
		return nil, response.ErrorMessage
	case protos.ErrorMessage_USER_NOT_EXISTS:
		return nil, response.ErrorMessage
	case protos.ErrorMessage_OK:
		//解析从CA得到的用户证书
		cert := &utils.Certificate{}
		err = json.Unmarshal(response.Content, cert)
		if err != nil {
			return nil, response.ErrorMessage
		}
		//获取CA的公钥
		request := &protos.GetCAPublicKeyRequest{}
		caResponse, err := c.GetCAPublicKey(ctx, request)
		if err != nil {
			return nil, response.ErrorMessage
		}
		caPubKey := &rsa.PublicKey{}
		err = json.Unmarshal(caResponse.Data, caPubKey)
		if err != nil {
			return nil, response.ErrorMessage
		}

		//用CA的公钥验证证书是否合法，这一步在客户端做而不是在KeyServer会更加安全,再验证一次用户名是为了防止有人冒充CA发送他自己的证书，这样可以绕过认证机制
		if !utils.VerifyCert(cert, caPubKey) || cert.Info.Username != username {
			return nil, response.ErrorMessage
		}

		return &cert.Info.PublicKey, response.ErrorMessage
	}
	return nil, protos.ErrorMessage_SERVER_ERROR
}
