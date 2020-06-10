package utils

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

/**
@params: bits为生成密钥的长度，默认采用1024位
@return: 生成密钥对文件
*/
func GenerateRSAKey(bits int, dir, username string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create(path.Join(dir, username+"private.pem"))
	if err != nil {
		return err
	}
	defer file.Close()
	err = pem.Encode(file, block)

	if err != nil {
		return err
	}

	publicKey := &privateKey.PublicKey
	derPKIX, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPKIX,
	}
	file, err = os.Create(path.Join(dir, username+"public.pem"))
	if err != nil {
		return err
	}

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

/**
@author ytw
公钥加密
@params:sourceData 源数据
		key 公钥
@returns:公钥加密后的string
*/
func PublickeyEncrypt(sourceData []byte, key *rsa.PublicKey) ([]byte, error) {
	res, err := rsa.EncryptPKCS1v15(rand.Reader, key, sourceData)
	if err != nil {
		return nil, err
	}
	return res, nil
}

/**
@author js
随机生成字符串，用来作为salt
*/
func GenerateRandomBytes(bits int) []byte {
	b := make([]byte, bits)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return b
}

/**
@author js
加盐哈希
*/
func HashWithSalt(passwd, salt string) string {
	h := sha256.New()
	h.Write([]byte(passwd + salt))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

/**
@author js
对文件进行散列，用于进行数字签名
*/
func HashFile(path string) ([]byte, error) {
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("读取文件出错 %v", err.Error())
		return nil, err
	}

	h := sha256.New()
	h.Write(buffer)
	return h.Sum(nil), nil
}

/**
@author js
从私钥文件读取私钥
*/
func GetPrivateKeyFromFile(filename string) *rsa.PrivateKey {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := FileToBytes(file)

	block, _ := pem.Decode(buf)

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return privateKey
}

/**
从公钥文件读取公钥
*/
func GetPublicKeyFromFile(filename string) *rsa.PublicKey {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	buf := FileToBytes(file)
	block, _ := pem.Decode(buf)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil
	}
	return publicKeyInterface.(*rsa.PublicKey)
}

/**
@author js
生成数字签名
@params:sourceData 源数据
		key 私钥
@returns:数字签名
*/
func GenerateDS(sourceData []byte, key *rsa.PrivateKey) ([]byte, error) {
	h := sha256.New()
	h.Write(sourceData)
	hashValue := h.Sum(nil)

	res, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashValue)
	if err != nil {
		return nil, err
	}
	return res, nil
}

/**
@author js
@params:signedData 数字签名
		sourceData 源数据
		key 公钥
*/
func VerifyDS(ds, sourceData []byte, key *rsa.PublicKey) bool {
	h := sha256.New()
	h.Write(sourceData)
	hashValue := h.Sum(nil)

	err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashValue, ds)
	return err == nil
}

/**
@author js
对文件用私钥进行签名，返回一个签名文件，后缀名为.sig
*/
func SignatureFile(file *os.File, key *rsa.PrivateKey) ([]byte, error) {
	buf := FileToBytes(file)
	ds, err := GenerateDS(buf, key)
	if err != nil {
		log.Printf("生成数字签名失败 %v", err.Error())
		return nil, err
	}

	return ds, nil
}

/**
@author js
用来生成AES密钥的结构
*/
type combineKey struct {
	secondKey  string
	privateKey rsa.PrivateKey
}

/**
@author js
使用二级密码和私钥生成对称密钥用于加密文件,返回的[]byte长度为16
@params secondKey 二级密码
		priKey 用户的私钥
@returns 二进制的对称密钥,error
*/
func GenerateSessionKeyWithSecondKey(secondKey string, priKey *rsa.PrivateKey) ([]byte, error) {
	comb := &combineKey{
		secondKey:  secondKey,
		privateKey: *priKey,
	}
	serializedComb, err := json.Marshal(comb)
	if err != nil {
		log.Printf("序列化密钥失败 %v", err.Error())
		return nil, err
	}

	h := md5.New()
	h.Write(serializedComb)
	hash := h.Sum(nil)
	return hash[0:16], nil
}

/**
@author ytw
对称密钥加密文件
*/
func AESEncryptFileToBytes(path string, key []byte) ([]byte, error) {
	file, _ := ReadFile(path)
	encrypt, err := aesEncrypt(file, key)
	if err != nil {
		log.Printf("文件加密失败 %v", err.Error())
		return nil, err
	}
	return encrypt, err
}

/**
@author ytw
对称密钥加密文件后生成文件
*/
func AESEncryptFileToFile(path string, key []byte, destPath string) error {
	fileToBytes, err := AESEncryptFileToBytes(path, key)
	if err != nil {
		log.Printf("文件加密失败 %v", err.Error())
		return err
	}
	err = WriteFile(destPath, fileToBytes)
	if err != nil {
		log.Printf("文件加密失败 %v", err.Error())
		return err
	}
	return nil
}

/**
@author js
AES解密到文件
*/
func AESDecryptToFile(data, key []byte, destPath string) error {
	data, err := aesDecrypt(data, key)
	if err != nil {
		log.Printf("文件解密失败 %v", err.Error())
		return err
	}
	err = WriteFile(destPath, data)
	if err != nil {
		log.Printf("文件解密失败 %v", err.Error())
		return err
	}
	return nil
}

/**
@author ytw
Aes填充
*/
func aesPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

/**
@author ytw
去除Aes填充
*/
func aesUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

/**
@author ytw
Aes加密
@params:origData 明文
		key 密钥
@returns:密文
*/
func aesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = aesPadding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

/**
@author ytw
Aes加密
@params:crypted 密文
		key 密钥
@returns:明文
*/
func aesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = aesUnPadding(origData)
	return origData, nil
}
