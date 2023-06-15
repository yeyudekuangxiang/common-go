package encrypttool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

func Md5(str string) string {
	encrypt := md5.New()
	encrypt.Write([]byte(str))
	md5Data := encrypt.Sum([]byte(""))
	return hex.EncodeToString(md5Data)
}

// Deprecated: 请使用 AesEncryptCBCPKCS7 代替
func AesEncrypt(orig string, key, iv string) string {
	return base64.StdEncoding.EncodeToString(AesEncryptCBCPKCS7([]byte(orig), []byte(key), []byte(iv)))
}
func AesEncryptCBCPKCS7(origData, key, iv []byte) []byte {
	// 分组秘钥 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
	block, _ := aes.NewCipher(key)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = pkcs7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return cryted
}
func AesDecryptCBCPKCS7(cryted, key, iv []byte) ([]byte, error) {
	// 分组秘钥
	block, _ := aes.NewCipher(key)
	// 加密模式 cbc
	blockMode := cipher.NewCBCDecrypter(block, iv)
	// 创建数组
	origin := make([]byte, len(cryted))
	// 解密
	blockMode.CryptBlocks(origin, cryted)
	// 去补全码
	origin, err := pkcs7Trimming(origin)
	if err != nil {
		return nil, err
	}
	return origin, nil
}

// Deprecated: 请使用 AesDecryptCBCPKCS7 代替
func AesDecrypt(cryted, key, iv string) (string, error) {
	// 转成字节数组
	crytedByte, err := base64.StdEncoding.DecodeString(cryted)
	if err != nil {
		return "", err
	}
	deData, err := AesDecryptCBCPKCS7(crytedByte, []byte(key), []byte(iv))
	if err != nil {
		return "", err
	}
	return string(deData), nil
}

func HMacMd5(orig, key string) string {
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(orig))
	return hex.EncodeToString(h.Sum([]byte("")))
}

// 生成 HMAC-MD5 签名
func hmacMD5Sign(key []byte, data []byte) string {
	h := hmac.New(md5.New, key)
	h.Write(data)
	signatureByte := h.Sum(nil)
	signature := hex.EncodeToString(signatureByte)
	return signature
}

// 补码
func pkcs7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去码
func pkcs7Trimming(encrypt []byte) ([]byte, error) {
	length := len(encrypt)
	if length == 0 {
		return nil, errors.New("加密字符串错误! ")
	}
	//获取填充的个数
	unPadding := int(encrypt[length-1])
	return encrypt[:(length - unPadding)], nil
}

func Sha1(str string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
}

func HMacSha1(orig, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(orig))
	return hex.EncodeToString(h.Sum([]byte("")))
}

func Sha256Byte(data []byte) string {
	encrypt := sha256.Sum256(data)
	return hex.EncodeToString(encrypt[:])
}
func Sha256(str string) string {
	return Sha256Byte([]byte(str))
}

//AES/ECB/PKCS7Padding **********开始 ***********

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func EcbDecrypt(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}
	return PKCS7UnPadding(decrypted)
}

func EcbEncryptCBC(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	data = PKCS7Padding(data, block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}
	return decrypted
}

//AES/ECB/PKCS7Padding **********结束 ***********

//AES/CBC/PKCS5Padding **********开始 ***********

//@brief:填充明文

func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

//@brief:去除填充数据

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//@brief:AES/PKCS5

func AesEncryptPKCS5(origData, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv) //初始向量的长度必须等于块block的长度16字节
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//@brief:AES解密

func AesDecryptPKCS5(crypted, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//AES分组长度为128位，所以blockSize=16，单位字节
	//blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv) //初始向量的长度必须等于块block的长度16字节
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}
