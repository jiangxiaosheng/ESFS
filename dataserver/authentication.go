package main

import (
	clicommon "ESFS2.0/client/common"
	datacommon "ESFS2.0/dataserver/common"
	"ESFS2.0/keyserver/common"
	"ESFS2.0/message/protos"
	"ESFS2.0/utils"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"path"
	"time"
)

/**
@author js
*/
func (s *dataServer) Login(ctx context.Context, req *protos.LoginRequest) (*protos.LoginResponse, error) {
	//获取数据库连接
	db, err := common.GetDBConnection()
	if err != nil {
		log.Printf("连接数据库失败 %v", err.Error())
		return &protos.LoginResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR}, err
	}

	//查询用户是否存在
	sql := fmt.Sprintf("select password_hash,salt from users where username='%s'", req.Username)
	res, err := common.DoQuery(sql, db)
	if err != nil {
		log.Printf("查询数据库失败 %v", err.Error())
		return &protos.LoginResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR}, err
	}

	var passwordHash, salt string
	exists := false
	if res.Next() {
		exists = true
		res.Scan(&passwordHash, &salt)
	}
	if exists == false {
		return &protos.LoginResponse{
			ErrorMessage: protos.ErrorMessage_USER_NOT_EXISTS}, nil
	}

	//检查密码是否正确
	user_passwordHash := utils.HashWithSalt(req.Password, salt)
	if user_passwordHash == passwordHash {
		return &protos.LoginResponse{
			ErrorMessage: protos.ErrorMessage_OK,
		}, nil
	}

	return &protos.LoginResponse{
		ErrorMessage: protos.ErrorMessage_PASSWORD_WRONG,
	}, nil
}

/**
@author js
*/
func (s *dataServer) Register(ctx context.Context, req *protos.RegisterRequest) (*protos.RegisterResponse, error) {
	//获取数据库连接
	db, err := common.GetDBConnection()
	if err != nil {
		log.Printf("连接数据库失败 %v", err.Error())
		return &protos.RegisterResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR}, err
	}

	//查询用户是否存在
	sql := fmt.Sprintf("select username from users where username='%s'", req.Username)
	res, err := common.DoQuery(sql, db)
	if err != nil {
		log.Printf("查询数据库失败 %v", err.Error())
		return &protos.RegisterResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR}, err
	}

	if res.Next() { //如果已存在，则返回失败
		return &protos.RegisterResponse{
			ErrorMessage: protos.ErrorMessage_USER_ALREADY_EXISTS,
		}, nil
	}

	//验证证书
	caPubKey := getCAPublicKey()
	if caPubKey == nil {
		return &protos.RegisterResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
		}, nil
	}
	cert := &utils.Certificate{}
	err = json.Unmarshal(req.CertData, cert)
	if err != nil {
		return &protos.RegisterResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
		}, nil
	}
	if !utils.VerifyCert(cert, caPubKey) { //验证不通过
		return &protos.RegisterResponse{
			ErrorMessage: protos.ErrorMessage_CERTIFICATE_INVALID,
		}, nil
	}

	username := req.Username
	password := req.Password
	defaultSecondKey := req.DefaultSecondKey
	salt := base64.URLEncoding.EncodeToString(utils.GenerateRandomBytes(32)) //随机生成salt
	passwordHash := utils.HashWithSalt(password, salt)
	sql = fmt.Sprintf("insert into users values('%s','%s','%s','%s')", username, passwordHash, salt, defaultSecondKey)

	//向数据库中插入数据
	_, err = common.DoExecTx(sql, db)
	if err != nil {
		log.Printf("数据库执行插入事务失败 %v", err.Error())
		return &protos.RegisterResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
		}, nil
	}

	//创建用户文件目录
	err = os.Mkdir(path.Join(common.BaseDir, "dataserver", "data", username), os.ModePerm)
	if err != nil {
		log.Printf("创建目录失败 %v", err.Error())
		return &protos.RegisterResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
		}, err
	}

	return &protos.RegisterResponse{
		ErrorMessage: protos.ErrorMessage_OK,
	}, nil
}

/**
@author ytw
加密结果存在access表中（可以为多个）
*/
func (s *dataServer) SaveSharedResult(ctx context.Context, req *protos.SaveSharedResultRequest) (*protos.SaveSharedResultResponse, error) {
	db, err := common.GetDBConnection()
	if err != nil {
		log.Printf("连接数据库失败 %v", err)
		return &protos.SaveSharedResultResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
		}, err
	}
	//查询用户、对应共享文件是否存在
	for i, filename := range req.Filenames {
		sql := fmt.Sprintf("select * from access where username='%s' and filename='%s'and authorized_user='%s'",
			req.Username, filename, req.AuthorizedUsername)
		res, err := common.DoQuery(sql, db)
		if err != nil {
			log.Printf("查询数据库失败 %v", err)
			return &protos.SaveSharedResultResponse{
				ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
			}, err
		}
		if res.Next() { //如果已存在
			sql := fmt.Sprintf("update access set share_key='%s'where username='%s' and filename='%s'and authorized_user='%s'and share_key<>'%s'",
				req.ShareKeys[i], req.Username, filename, req.AuthorizedUsername, req.ShareKeys[i])
			fmt.Printf("***update***sql:" + sql + "\n")
			_, err := common.DoExecTx(sql, db)
			if err != nil {
				log.Printf("更新数据库失败 %v", err)
				return &protos.SaveSharedResultResponse{
					ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
				}, err
			}
		} else { //不存在
			sql = fmt.Sprintf("insert into access values('%s','%s','%s','%s')",
				req.Username, filename, req.AuthorizedUsername, req.ShareKeys[i])
			fmt.Printf("***insert***sql:" + sql + "\n")
			_, err := common.DoExecTx(sql, db)
			if err != nil {
				log.Printf("插入数据库失败 %v", err)
				return &protos.SaveSharedResultResponse{
					ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
				}, err
			}
		}
	}
	return &protos.SaveSharedResultResponse{
		ErrorMessage: protos.ErrorMessage_OK,
	}, nil
}

/**
@author yyx
获取某个用户指定文件的二级密码（可以为多个）
*/
func (s *dataServer) GetSecondKey(ctx context.Context, req *protos.GetSecondKeyRequest) (*protos.GetSecondKeyResponse, error) {
	db, err := common.GetDBConnection()
	if err != nil {
		log.Printf("连接数据库失败 %v", err)
		return &protos.GetSecondKeyResponse{
			ErrorMessage:      protos.ErrorMessage_SERVER_ERROR,
			SecondKeysMapData: nil,
		}, err
	}
	//查询用户是否存在
	sql := fmt.Sprintf("select * from users where username='%s'", req.Username)
	res, err := common.DoQuery(sql, db)
	if err != nil {
		log.Printf("查询数据库失败 %v", err)
		return &protos.GetSecondKeyResponse{
			ErrorMessage:      protos.ErrorMessage_SERVER_ERROR,
			SecondKeysMapData: nil,
		}, err
	}

	if !res.Next() {
		return &protos.GetSecondKeyResponse{
			ErrorMessage:      protos.ErrorMessage_USER_NOT_EXISTS,
			SecondKeysMapData: nil,
		}, err
	}

	ranges := ""
	for i, name := range req.Filenames {
		ranges += fmt.Sprintf("'%s'", name)
		if i != len(req.Filenames)-1 {
			ranges += ","
		}
	}

	sql = fmt.Sprintf("select filename,secondKey from metadata where username='%s' and filename in (%s)", req.Username, ranges)
	res, err = datacommon.DoQuery(sql, db)
	if err != nil {
		log.Printf("查询数据库失败 %v", err)
		return &protos.GetSecondKeyResponse{
			ErrorMessage:      protos.ErrorMessage_SERVER_ERROR,
			SecondKeysMapData: nil,
		}, err
	}

	m := make(map[string]string)
	var filename, secondkey string
	for res.Next() {
		res.Scan(&filename, &secondkey)
		m[filename] = secondkey
	}

	serializedData, err := json.Marshal(m)
	if err != nil {
		log.Printf("序列化失败 %v", err)
		return &protos.GetSecondKeyResponse{
			ErrorMessage:      protos.ErrorMessage_SERVER_ERROR,
			SecondKeysMapData: nil,
		}, err
	}

	return &protos.GetSecondKeyResponse{
		ErrorMessage:      protos.ErrorMessage_OK,
		SecondKeysMapData: serializedData,
	}, nil
}

/**
@author js
查询用户是否存在
*/
func (s *dataServer) CheckUserExits(ctx context.Context, req *protos.CheckUserExistsRequest) (*protos.CheckUserExistsResponse, error) {
	db, err := common.GetDBConnection()
	if err != nil {
		log.Printf("连接数据库失败 %v", err.Error())
		return &protos.CheckUserExistsResponse{
			Exists:       false,
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
		}, err
	}

	sql := fmt.Sprintf("select * from users where username='%s'", req.Username)
	rows, err := common.DoQuery(sql, db)
	if err != nil {
		log.Printf("查询数据库失败 %v", err.Error())
		return &protos.CheckUserExistsResponse{
			Exists:       false,
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
		}, err
	}
	if rows.Next() {
		return &protos.CheckUserExistsResponse{
			Exists:       true,
			ErrorMessage: protos.ErrorMessage_OK,
		}, nil
	}
	return &protos.CheckUserExistsResponse{
		Exists:       false,
		ErrorMessage: protos.ErrorMessage_OK,
	}, nil
}

/**
@author js
*/
func getCAPublicKey() *rsa.PublicKey {
	c, conn, err := clicommon.GetCAClient()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := &protos.GetCAPublicKeyRequest{}
	response, err := c.GetCAPublicKey(ctx, request)
	if err != nil {
		return nil
	}
	caPubKey := &rsa.PublicKey{}
	err = json.Unmarshal(response.Data, caPubKey)
	if err != nil {
		return nil
	}
	return caPubKey
}

/**
@author js
获取指定用户的默认二级密码
*/
func getDefaultSecondKey(username string) (string, error) {
	db, err := common.GetDBConnection()
	if err != nil {
		log.Printf("获取数据库连接失败 %v", err.Error())
		return "", err
	}

	sql := fmt.Sprintf("select defaultSecondKey from users where username='%s'", username)
	res, err := common.DoQuery(sql, db)
	if err != nil {
		log.Printf("查询数据库失败 %v", err.Error())
		return "", err
	}

	if res.Next() {
		var secondKey string
		res.Scan(&secondKey)
		return secondKey, nil
	} else {
		return "", errors.New("用户不存在")
	}
}
