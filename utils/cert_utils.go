package utils

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
)

type Certificate struct {
	Info CertInfo
	Sig  string //byte数组base64编码
}

type CertInfo struct {
	Username  string
	PublicKey rsa.PublicKey
}

/**
@author js
*/
func GenerateCertToFile(username string, pubkey *rsa.PublicKey, priKey *rsa.PrivateKey, dir string) error {
	certInfo := CertInfo{
		Username:  username,
		PublicKey: *pubkey,
	}

	serializedKey, err := json.Marshal(certInfo) //序列化用户信息
	if err != nil {
		log.Printf("序列化失败 %v", err.Error())
		return err
	}
	sig, err := GenerateDS(serializedKey, priKey) //生成数字签名
	if err != nil {
		log.Printf("生成数字签名失败 %v", err.Error())
		return err
	}
	cert := Certificate{
		Info: certInfo,
		Sig:  base64.URLEncoding.EncodeToString(sig),
	}

	certFile, err := os.Create(dir)
	if err != nil {
		log.Printf("证书输出文件创建失败 %v", err.Error())
		return err
	}

	serializedCert, err := json.Marshal(cert) //序列化Certificate
	if err != nil {
		log.Printf("序列化失败 %v", err.Error())
		return err
	}

	//var headers = make(map[string]string)
	//headers["username"] = "memeshe"
	//headers["publickey"] = base64.URLEncoding.EncodeToString(pubkey.N.Bytes())
	err = pem.Encode(certFile, &pem.Block{
		Type: "CERTIFICATE",
		//Headers: headers,
		Bytes: serializedCert,
	})

	certFile.Close()
	return nil
}

/**
@author js
*/
func GenerateCertToBytes(username string, pubKey *rsa.PublicKey, priKey *rsa.PrivateKey) ([]byte, error) {
	certInfo := CertInfo{
		Username:  username,
		PublicKey: *pubKey,
	}

	serializedKey, err := json.Marshal(certInfo) //序列化用户信息
	if err != nil {
		log.Printf("序列化失败 %v", err.Error())
		return nil, err
	}
	sig, err := GenerateDS(serializedKey, priKey) //生成数字签名
	if err != nil {
		log.Printf("生成数字签名失败 %v", err.Error())
		return nil, err
	}
	cert := Certificate{
		Info: certInfo,
		Sig:  base64.URLEncoding.EncodeToString(sig),
	}

	serializedCert, err := json.Marshal(cert) //序列化Certificate
	if err != nil {
		log.Printf("序列化失败 %v", err.Error())
		return nil, err
	}

	//var headers = make(map[string]string)
	//headers["username"] = username
	//headers["publickey"] = base64.URLEncoding.EncodeToString(pubKey.N.Bytes())
	block := &pem.Block{
		Type: "CERTIFICATE",
		//Headers: headers,
		Bytes: serializedCert,
	}
	b := pem.EncodeToMemory(block)
	return b, nil
}

/**
从证书文件读取Certificate
*/
func ReadCertFromFile(path string) (*Certificate, error) {
	encodedCert, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("证书格式错误 %v", err.Error())
		return nil, err
	}
	block, _ := pem.Decode(encodedCert)

	cert := &Certificate{}
	err = json.Unmarshal(block.Bytes, cert)
	if err != nil {
		log.Printf("序列化失败 %v", err.Error())
		return nil, err
	}
	return cert, nil
}

/**
@author js
*/
func ReadCertFromBytes(source []byte) (*Certificate, error) {
	block, _ := pem.Decode(source)

	cert := &Certificate{}
	err := json.Unmarshal(block.Bytes, cert)
	if err != nil {
		log.Printf("反序列化失败 %v", err.Error())
		return nil, err
	}
	return cert, nil
}

/**
@author js
验证证书
*/
func VerifyCert(cert *Certificate, pubKey *rsa.PublicKey) bool {
	sig, err := base64.URLEncoding.DecodeString(cert.Sig)
	if err != nil {
		return false
	}
	serializedData, err := json.Marshal(cert.Info)
	if err != nil {
		return false
	}
	return VerifyDS(sig, serializedData, pubKey)
}
