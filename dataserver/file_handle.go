package main

import (
	datacommon "ESFS2.0/dataserver/common"
	"ESFS2.0/keyserver/common"
	"ESFS2.0/message"
	"ESFS2.0/message/protos"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
)

func fileSocketServer() {
	port := 8959
	host := "0.0.0.0"
	addr := fmt.Sprintf("%s:%d", host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[socket] failed to listen: %v", err.Error())
	}
	defer lis.Close()

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("建立socket连接失败 %v", err.Error())
			continue
		}

		go func(conn net.Conn) {
			msg := &message.FileSocketMessage{}
			buffer := make([]byte, 2048)
			n, err := conn.Read(buffer)
			err = json.Unmarshal(buffer[:n], msg)
			if err != nil {
				log.Printf("反序列化失败 %v", err.Error())
				return
			}

			conn.Write([]byte{1}) //给客户端一个信息，防止粘包

			if msg.Type == message.FILE_UPLOAD {
				file, err := os.Create(path.Join(common.BaseDir, "dataserver", "data", msg.UserName, msg.FileName[0]))
				if err != nil {
					log.Printf("打开文件失败 %v", err.Error())
					return
				}

				for {
					n, err := conn.Read(buffer)
					if err == io.EOF {
						break
					}
					_, err = file.Write(buffer[:n])
					if err != nil {
						log.Printf("写入文件失败 %v", err.Error())
						break
					}
				}
				file.Close()
			} else if msg.Type == message.FILE_DOWNLOAD {
				//TODO
				conn.Read(buffer)

				filenames := msg.FileName
				for _, filename := range filenames {
					file, err := os.Open(path.Join(common.BaseDir, "dataserver", "data", msg.UserName, filename))
					if err != nil {
						log.Printf("打开文件失败 %v", err.Error())
						return
					}

					for {
						n, err := file.Read(buffer)
						if err == io.EOF {
							signal := &message.SignalOver{}
							serializedData, _ := json.Marshal(signal)
							_, err = conn.Write(serializedData)
							if err != nil {
								log.Printf("写入socket失败 %v", err.Error())
							}
							break
						}

						_, err = conn.Write(buffer[:n])
						if err != nil {
							log.Printf("写入socket失败 %v", err.Error())
							break
						}
						conn.Read(buffer)
					}

					file.Close()
					conn.Read(buffer) //防止粘包

					sigFile, err := os.Open(path.Join(common.BaseDir, "dataserver", "data", msg.UserName, "."+filename+".sig"))
					if err != nil {
						log.Printf("打开签名文件失败 %v", err.Error())
						return
					}
					log.Printf("发送签名文件 %s", sigFile.Name())
					for {
						n, err := sigFile.Read(buffer)
						if err == io.EOF {
							signal := &message.SignalOver{}
							serializedData, _ := json.Marshal(signal)
							_, err = conn.Write(serializedData)
							if err != nil {
								log.Printf("写入socket失败 %v", err.Error())
							}
							break
						}

						_, err = conn.Write(buffer[:n])
						if err != nil {
							log.Printf("写入socket失败 %v", err.Error())
							break
						}
						conn.Read(buffer)
					}
					log.Printf("签名文件发送完毕 %s", sigFile.Name())
					sigFile.Close()
					conn.Read(buffer)
				}
			}
		}(conn)
	}
}

/**
@author js
*/
func (s *dataServer) UploadDS(ctx context.Context, req *protos.UploadDSRequest) (*protos.UploadDSResponse, error) {
	file, err := os.Create(path.Join(common.BaseDir, "dataserver", "data", req.Username, "."+req.Filename+".sig"))
	if err != nil {
		log.Printf("创建文件失败 %v", err.Error())
		return &protos.UploadDSResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
		}, err
	}

	err = ioutil.WriteFile(file.Name(), req.DsData, 0644)
	if err != nil {
		log.Printf("写入签名数据失败 %v", err.Error())
		return &protos.UploadDSResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
		}, err
	}
	file.Close()

	return &protos.UploadDSResponse{ErrorMessage: protos.ErrorMessage_OK}, nil
}

/**
@author js
*/
func (s *dataServer) UploadPrepare(ctx context.Context, req *protos.UploadPrepareRequest) (*protos.UploadPrepareResponse, error) {
	//反序列化，获取文件信息
	fileInfo := &message.FileInfo{}
	err := json.Unmarshal(req.FileInfo, fileInfo)
	if err != nil {
		log.Printf("反序列化失败 %v", err.Error())
		return &protos.UploadPrepareResponse{
			ErrorMessage:     protos.ErrorMessage_SERVER_ERROR,
			DefaultSecondKey: "",
		}, err
	}

	//创建指定文件
	//Create函数若文件已存在则会截断，不存在则新建
	file, err := os.Create(path.Join(common.BaseDir, "dataserver", "data", req.Username, fileInfo.Name))
	fmt.Printf(fileInfo.Name)
	if err != nil {
		log.Printf("创建文件失败 %v", err.Error())
		return &protos.UploadPrepareResponse{
			ErrorMessage:     protos.ErrorMessage_SERVER_ERROR,
			DefaultSecondKey: "",
		}, err
	}
	defer file.Close()

	//从数据库中读出用户的默认二级密码
	secondKey, err := getDefaultSecondKey(req.Username)
	if err != nil {
		log.Printf("获取二级密码失败 %v", err.Error())
		return &protos.UploadPrepareResponse{
			ErrorMessage:     protos.ErrorMessage_SERVER_ERROR,
			DefaultSecondKey: "",
		}, err
	}

	if req.SecondKey != "" { //指定新二级密码
		err = updateMetadata(req.Username, fileInfo.Name, req.SecondKey)
	} else { //使用默认二级密码
		err = updateMetadata(req.Username, fileInfo.Name, secondKey)
	}

	if err != nil {
		log.Printf("数据库更新失败 %v", err.Error())
	}

	return &protos.UploadPrepareResponse{
		ErrorMessage:     protos.ErrorMessage_OK,
		DefaultSecondKey: secondKey,
	}, nil
}

func (s *dataServer) DownloadPrepare(ctx context.Context, req *protos.DownloadPrepareRequest) (*protos.DownloadPrepareResponse, error) {
	db, err := datacommon.GetDBConnection()
	if err != nil {
		log.Printf("建立数据库连接失败 %v", err.Error())
		return &protos.DownloadPrepareResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
			Content:      nil,
		}, err
	}
	sql := fmt.Sprintf("select filename,secondKey from metadata where username='%s'", req.Username)
	res, err := datacommon.DoQuery(sql, db)
	if err != nil {
		log.Printf("查询metadata数据库失败 %v", err.Error())
		return &protos.DownloadPrepareResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
			Content:      nil,
		}, err
	}
	var filename, key string
	m := make(map[string]string)
	for res.Next() {
		res.Scan(&filename, &key)
		m[filename] = key
	}
	data, err := json.Marshal(m)
	if err != nil {
		log.Printf("序列化二级密码映射失败 %v", err.Error())
		return &protos.DownloadPrepareResponse{
			ErrorMessage: protos.ErrorMessage_SERVER_ERROR,
			Content:      nil,
		}, err
	}

	response := &protos.DownloadPrepareResponse{
		ErrorMessage: protos.ErrorMessage_OK,
		Content:      data,
	}
	return response, nil
}

func (s *dataServer) ListFiles(ctx context.Context, req *protos.ListFilesRequest) (*protos.ListFilesResponse, error) {
	fileDir := path.Join(common.BaseDir, "dataserver", "data", req.Username)
	files, err := ioutil.ReadDir(fileDir)
	if err != nil {
		log.Printf("读目录失败 %v", err.Error())
		return &protos.ListFilesResponse{
			Ok:       false,
			FileInfo: nil,
		}, err
	}

	var filesArray [][]byte
	for _, f := range files {
		if path.Ext(f.Name()) == ".sig" && string(f.Name()[0]) == "." {
			continue
		}

		tmp := &message.FileInfo{
			Name:    f.Name(),
			Mode:    f.Mode(),
			Size:    f.Size(),
			ModTime: f.ModTime(),
		}
		serializedData, err := json.Marshal(tmp)
		if err != nil {
			log.Printf("反序列化文件信息失败 %v", err.Error())
			continue
		}

		filesArray = append(filesArray, serializedData)
	}

	return &protos.ListFilesResponse{
		Ok:       true,
		FileInfo: filesArray,
	}, nil
}

func (s *dataServer) RemoveFiles(ctx context.Context, req *protos.RemoveFilesRequest) (*protos.RemoveFilesResponse, error) {
	for _, filename := range req.Filenames {
		err := os.Remove(path.Join(datacommon.BaseDir, "dataserver", "data", req.Username, filename))           //删除原文件
		err = os.Remove(path.Join(datacommon.BaseDir, "dataserver", "data", req.Username, "."+filename+".sig")) // 删除签名文件
		if err != nil {
			log.Printf("删除文件失败 %v", err.Error())
			return &protos.RemoveFilesResponse{ErrorMessage: protos.ErrorMessage_SERVER_ERROR}, err
		}
	}
	return &protos.RemoveFilesResponse{ErrorMessage: protos.ErrorMessage_OK}, nil
}

func updateMetadata(username, filename, secondKey string) error {
	db, err := datacommon.GetDBConnection()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("select * from metadata where username='%s' and filename='%s'", username, filename)
	rows, err := datacommon.DoQuery(sql, db)
	if err != nil {
		return err
	}
	if rows.Next() {
		sql = fmt.Sprintf("update metadata set secondKey='%s' where username='%s' and filename='%s'", secondKey, username, filename)
	} else {
		sql = fmt.Sprintf("insert into metadata values('%s','%s','%s')", username, filename, secondKey)
	}
	_, err = datacommon.DoExecTx(sql, db)
	if err != nil {
		return err
	}
	return nil
}
