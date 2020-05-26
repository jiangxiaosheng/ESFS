package keyserver

import (
	"ESFS2.0/keyserver/common"
	"ESFS2.0/message/protos"
	"context"
	"log"
)

type server struct {
	protos.UnimplementedKeyStoreServer
}

func (s *server) GetCert(ctx context.Context, req *protos.GetCertRequest) (*protos.GetCertResponse, error) {
	db, err := common.GetDBConnection()
	if err != nil {
		log.Println(err.Error())
		return &protos.GetCertResponse{
			Ok:      false,
			Content: nil}, err
	}

	username := req.Username
	table := utils.GetConfig("common.keyserver.table")

	sql := fmt.Sprintf("select cert from %s where username='%s'", table, username)
	rows, err := utils.DoQuery(sql, db)

	if err != nil {
		fmt.Println("查询证书失败", err.Error())
		return &proto.Cert{Ok: false, Content: nil}, err
	}

	if rows.Next() { //如果存在相应的证书
		content := new([]byte)
		rows.Scan(content)

		cert := &proto.Cert{Ok: true, Content: *content}
		return cert, nil
	} else {
		cert := &proto.Cert{
			Ok:      false,
			Content: nil,
		}
		return cert, nil
	}
}
