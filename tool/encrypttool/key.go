package encrypttool

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// LoadRsaPrivateKeyPKCS8Pem 通过私钥的文本内容加载rsa私钥
func LoadRsaPrivateKeyPKCS8Pem(privateKeyStr string) (privateKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, fmt.Errorf("decode private key err")
	}
	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("the kind of PEM should be PRVATE KEY")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key err:%s", err.Error())
	}
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not a RSA private key")
	}
	return privateKey, nil
}

// LoadRsaPublicKeyPem 通过公钥的文本内容加载rsa公钥
func LoadRsaPublicKeyPem(publicKeyStr string) (publicKey *rsa.PublicKey, err error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return nil, errors.New("decode public key error")
	}
	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("the kind of PEM should be PUBLIC KEY")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key err:%s", err.Error())
	}
	publicKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("%s is not rsa public key", publicKeyStr)
	}
	return publicKey, nil
}

// LoadRsaPrivateKeyPKCS8Base64 通过私钥的base64字符串加载rsa私钥
func LoadRsaPrivateKeyPKCS8Base64(privateKeyStr string) (privateKey *rsa.PrivateKey, err error) {
	data, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("decode private key err:%s", err.Error())
	}
	key, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("parse private key err:%s", err.Error())
	}
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not a RSA private key")
	}
	return privateKey, nil
}

// LoadRsaPublicKeyPKCS8Base64 通过公钥钥的base64字符串加载rsa私钥
func LoadRsaPublicKeyPKCS8Base64(publicKeyStr string) (privateKey *rsa.PublicKey, err error) {
	data, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return nil, fmt.Errorf("decode public key err:%s", err.Error())
	}
	key, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, fmt.Errorf("parse public key err:%s", err.Error())
	}
	publicKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not a RSA public key")
	}
	return publicKey, nil
}
