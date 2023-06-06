package ncount

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func httpPost(url string, form url.Values) ([]byte, error) {
	resp, err := http.PostForm(url, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// 1. 用新账通平台公钥对json字符串进行非对称加密；
// 2. 对加密后的二进制转 Base64 编码
// 通过单元测试
func Encrpt(message []byte, key string) ([]byte, error) {
	publicKeyBlock, _ := pem.Decode([]byte(key))
	if publicKeyBlock == nil {
		return nil, fmt.Errorf("Key is invalid")
	}
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, errors.New("Encrpt ParsePKIXPublicKey :" + err.Error())
	}
	cipherByte, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), message)
	if err != nil {
		return nil, err
	}
	encoded := base64.StdEncoding.EncodeToString(cipherByte)
	return []byte(encoded), nil
}

// 完成签名认证 ，通过单元测试
func Sign(message []byte, privateKeyString string) (string, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyString)
	if err != nil {
		return "", err
	}
	// 将私钥解析为 *rsa.PrivateKey 对象
	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return "", err
	}
	hashed := sha1.Sum(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// 非对称加密-分段
func RsaEncryptBlock(message []byte, key string) (bytesEncrypt []byte, err error) {
	publicKeyBlock, _ := pem.Decode([]byte(key))
	if publicKeyBlock == nil {
		return nil, errors.New("Key is invalid")
	}
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, errors.New("Encrpt ParsePKIXPublicKey :" + err.Error())
	}
	keySize, srcSize := publicKey.(*rsa.PublicKey).Size(), len(message)
	pub := publicKey.(*rsa.PublicKey)
	//单次加密的长度需要减掉padding的长度，PKCS1为11
	offSet, once := 0, keySize-11
	buffer := bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + once
		if endIndex > srcSize {
			endIndex = srcSize
		}
		// 加密一部分
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, pub, message[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	bytesEncrypt = buffer.Bytes()
	encoded := base64.StdEncoding.EncodeToString(bytesEncrypt)
	return []byte(encoded), nil
}
