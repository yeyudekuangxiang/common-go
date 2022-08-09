package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type orig struct {
	Total             int               `json:"total"`
	StationStatusInfo stationStatusInfo `json:"stationStatusInfo"`
}

type stationStatusInfo struct {
	OperationID          string               `json:"operationID"`
	StationID            string               `json:"stationID"`
	ConnectorStatusInfos connectorStatusInfos `json:"connectorStatusInfos"`
}

type connectorStatusInfos struct {
	ConnectorID string `json:"connectorID"`
	EquipmentID string `json:"equipmentID"`
	Status      int    `json:"status"`
	CurrentA    int    `json:"currentA"`
	CurrentB    int    `json:"currentB"`
	CurrentC    int    `json:"currentC"`
	VoltageA    int    `json:"voltageA"`
	VoltageB    int    `json:"voltageB"`
	VoltageC    int    `json:"voltageC"`
	Soc         int    `json:"soc"`
}

func TestMd5(t *testing.T) {
	str := fmt.Sprintf("%x", md5.Sum([]byte("123456789ABCabc")))
	fmt.Println(str)
}

func TestAES(t *testing.T) {
	//key长度必须为16为
	origin := orig{
		Total: 1,
		StationStatusInfo: stationStatusInfo{
			OperationID: "123456789",
			StationID:   "1111111111111111",
			ConnectorStatusInfos: connectorStatusInfos{
				ConnectorID: "1",
				EquipmentID: "1000000000000000000000001",
				Status:      0,
				CurrentA:    0,
				CurrentB:    0,
				CurrentC:    0,
				VoltageA:    0,
				VoltageB:    0,
				VoltageC:    0,
				Soc:         10,
			},
		},
	}
	marshal, err := json.Marshal(origin)
	if err != nil {
		return
	}
	fmt.Println("原文：", string(marshal))

	key := "1234567890abcdef"
	encryptCode := aesEncrypt(string(marshal), key)
	fmt.Println("密文：", encryptCode)

	decryptCode := aesDecrypt(encryptCode, key)
	fmt.Println("解密结果：", decryptCode)

}

func aesEncrypt(orig, key string) string {
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

func aesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式 cbc
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	origin := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(origin, crytedByte)
	// 去补全码
	origin, err := PKCS7Trimming(origin)
	if err != nil {
		return err.Error()
	}
	return string(origin)
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

//去码
func PKCS7Trimming(encrypt []byte) ([]byte, error) {
	length := len(encrypt)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(encrypt[length-1])
	return encrypt[:(length - unPadding)], nil
}
func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
