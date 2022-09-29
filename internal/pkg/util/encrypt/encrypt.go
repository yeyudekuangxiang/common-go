package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func Md5(str string) string {
	return Md5Byte([]byte(str))
}
func Md5Byte(data []byte) string {
	encrypt := md5.New()
	encrypt.Write(data)
	md5Data := encrypt.Sum([]byte(""))
	return hex.EncodeToString(md5Data)
}
func AesEncrypt(orig string, key, iv string) string {
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = pkcs7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

func AesDecrypt(cryted, key, iv string) (string, error) {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)

	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	//blockSize := block.BlockSize()
	// 加密模式 cbc
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	// 创建数组
	origin := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(origin, crytedByte)
	// 去补全码
	origin, err := pkcs7Trimming(origin)
	if err != nil {
		return "", err
	}
	return string(origin), nil
}

func HMacMd5(orig, key string) string {
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(orig))
	return hex.EncodeToString(h.Sum([]byte("")))
}

//补码
func pkcs7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func pkcs7Trimming(encrypt []byte) ([]byte, error) {
	length := len(encrypt)
	if length == 0 {
		return nil, errors.New("加密字符串错误! ")
	}
	//获取填充的个数
	unPadding := int(encrypt[length-1])
	return encrypt[:(length - unPadding)], nil
}

func Sha256Byte(data []byte) string {
	encrypt := sha256.Sum256(data)
	return hex.EncodeToString(encrypt[:])
}
func Sha256(str string) string {
	return Sha256Byte([]byte(str))
}