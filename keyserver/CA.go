package main

import (
	"ESFS2.0/keyserver/common"
	"ESFS2.0/message/protos"
	"ESFS2.0/utils"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"path"
)

func (s *keyServer) GetCert(ctx context.Context, req *protos.GetCertRequest) (*protos.GetCertResponse, error) {
	db, err := common.GetDBConnection()
	if err != nil {
		log.Printf("获取数据库连接失败 %v", err.Error())
		return &protos.GetCertResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
			Content:      nil,
		}, err
	}

	username := req.Username

	sql := fmt.Sprintf("select cert from cert where username='%s'", username)
	rows, err := common.DoQuery(sql, db)

	if err != nil {
		fmt.Println("查询证书失败", err.Error())
		return &protos.GetCertResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
			Content:      nil,
		}, err
	}

	if rows.Next() { //如果存在相应的证书
		content := new([]byte)
		rows.Scan(content)

		cert := &protos.GetCertResponse{
			Content: *content,
		}
		return cert, nil
	} else {
		cert := &protos.GetCertResponse{
			Content: nil,
		}
		return cert, nil
	}
}

func (s *keyServer) SetCert(ctx context.Context, req *protos.SetCertRequest) (*protos.SetCertResponse, error) {
	db, err := common.GetDBConnection()
	if err != nil {
		log.Printf("建立数据库连接失败 %v", err.Error())
		return &protos.SetCertResponse{ErrorMessage: protos.ErrorMessage_SERVER_ERROR}, err
	}
	sql := fmt.Sprintf("select * from cert where username='%s'", req.Username)
	rows, err := common.DoQuery(sql, db)
	if rows.Next() {
		return &protos.SetCertResponse{ErrorMessage: protos.ErrorMessage_USER_ALREADY_EXISTS}, nil
	}

	priKey := getKeyServerPrivateKey()
	if priKey == nil {
		return &protos.SetCertResponse{ErrorMessage: protos.ErrorMessage_SERVER_ERROR}, err
	}
	userPubKey := &rsa.PublicKey{}
	err = json.Unmarshal(req.Content, userPubKey)
	if err != nil {
		return &protos.SetCertResponse{ErrorMessage: protos.ErrorMessage_SERVER_ERROR}, err
	}

	cert, err := utils.GenerateCertToBytes(req.Username, userPubKey, priKey)
	if err != nil {
		return &protos.SetCertResponse{ErrorMessage: protos.ErrorMessage_SERVER_ERROR}, err
	}

	sql = fmt.Sprintf("insert into cert values('%s','%s')", req.Username, cert)
	_, err = common.DoExecTx(sql, db)
	if err != nil {
		log.Printf("数据库事务执行失败 %v", err.Error())
		return &protos.SetCertResponse{ErrorMessage: protos.ErrorMessage_SERVER_ERROR}, nil
	}
	return &protos.SetCertResponse{ErrorMessage: protos.ErrorMessage_OK, CertData: cert}, nil
}

func (s *keyServer) GetCAPublicKey(ctx context.Context, req *protos.GetCAPublicKeyRequest) (*protos.GetCAPublicKeyResponse, error) {
	serializedData, err := json.Marshal(getKeyServerPublicKey())
	if err != nil {
		return &protos.GetCAPublicKeyResponse{
			Data: nil,
		}, err
	}
	return &protos.GetCAPublicKeyResponse{
		Data: serializedData,
	}, nil
}

func getKeyServerPublicKey() *rsa.PublicKey {
	return utils.GetPublicKeyFromFile(path.Join(common.BaseDir, "keyserver", "key", "public.pem"))
}

func getKeyServerPrivateKey() *rsa.PrivateKey {
	return utils.GetPrivateKeyFromFile(path.Join(common.BaseDir, "keyserver", "key", "private.pem"))
}
