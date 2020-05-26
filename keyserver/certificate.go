package keyserver

import (
	"ESFS2.0/keyserver/common"
	"ESFS2.0/message/protos"
	"context"
	"fmt"
	"log"
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
			Ok:      true,
			Content: *content,
		}
		return cert, nil
	} else {
		cert := &protos.GetCertResponse{
			Ok:      false,
			Content: nil,
		}
		return cert, nil
	}
}
