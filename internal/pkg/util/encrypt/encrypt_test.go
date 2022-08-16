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
	"encoding/json"
	"fmt"
	"strings"
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

func TestHMAC_Md5(t *testing.T) {
	key := "123456789"
	data := "l7B0BSEjFdzpyKzfOFpvg/SelCP802RItKYFPfSLRxJ3jf0bV19hvYOEktPAYW2nd7S8MBcyHYyacHKbISq5iTmDzG+ivnR+SZJv3USNTYVMz9rCQVSxd0cLlgsJauko79NnwOJbzDTyLooYolwz75gBOH2/xOMirpeEqRJrF/EQjWekJmGk9RtboXePu2rka+Xm51syBPhiXJAq0GfbfaFu9Nqs/e2Via/tEIM0lqvxfXQ6da6HrThsm5id4ClZFliOacRfrsPLRixS/IQYtksxghvJwbqOsbIsITail9Ayy4tKcogeEZiOO+4Ed264NSKmk7l3wKwJLAFjCFogBx8GE3OBz4pqcAn/ydA="
	time := "20160729142400"
	seq := "0001"
	secret := key + data + time + seq
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(secret))
	str := fmt.Sprintf("%x", hex.EncodeToString(h.Sum([]byte(""))))
	fmt.Printf("%d-%s\n", len(str), str)
	fmt.Println(len("745166E8C43C84D37FFECOF529C4136F"))
}

const (
	ipad = 0x36
	opad = 0x5c
)

var (
	d       byte
	istr    string
	dsSlice []byte
)

func TestHMAC_Md52(t *testing.T) {
	operatorID := "123456789"
	data := "l7B0BSEjFdzpyKzfOFpvg/SelCP802RItKYFPfSLRxJ3jf0bV19hvYOEktPAYW2nd7S8MBcyHYyacHKbISq5iTmDzG+ivnR+SZJv3USNTYVMz9rCQVSxd0cLlgsJauko79NnwOJbzDTyLooYolwz75gBOH2/xOMirpeEqRJrF/EQjWekJmGk9RtboXePu2rka+Xm51syBPhiXJAq0GfbfaFu9Nqs/e2Via/tEIM0lqvxfXQ6da6HrThsm5id4ClZFliOacRfrsPLRixS/IQYtksxghvJwbqOsbIsITail9Ayy4tKcogeEZiOO+4Ed264NSKmk7l3wKwJLAFjCFogBx8GE3OBz4pqcAn/ydA="
	timeStamp := "20160729142400"
	seq := "0001"
	sigSecret := "1234567890abcde"
	padding := 64 - len([]byte(sigSecret))%64
	strNo1 := fmt.Sprintf(sigSecret+"%s", bytes.Repeat([]byte("0"), padding))
	fmt.Printf("%d-%s\n", len(strNo1), strNo1)
	istr = xorEncode(strNo1, ipad)
	fmt.Printf("ipad明文：%s ipad密文：%s\n", strNo1, istr)
	//ipadDecode := xorDecode(dsSlice, istr, ipad)
	//fmt.Printf("ipad密文：%s ipad明文：%s\n", istr, ipadDecode)
	strMd5 := Md5(istr + operatorID + data + timeStamp + seq)
	fmt.Printf("md5密文：%s md5明文：%s\n", strMd5, "ipadDecode")
	ostr := xorEncode(strNo1, opad)
	fmt.Printf("opad明文：%s opad密文：%s\n", strNo1, ostr)
	strLast := strings.ToUpper(Md5(ostr + strMd5))
	fmt.Printf("last明文：%s last密文：%s\n", ostr+strMd5, strLast)
}

func xorEncode(str string, pad int32) string {
	for _, v := range str {
		d = byte(v ^ pad)
		dsSlice = append(dsSlice, d)
	}
	return string(dsSlice)
}
func xorDecode(secret []byte, str string, pad int32) string {
	secret = secret[:0]
	for _, v := range str {
		d = byte(v ^ pad)
		secret = append(secret, d)
	}
	return string(secret)
}

func TestHMAC_Sha256(t *testing.T) {
	key := "123456789"
	data := "l7B0BSEjFdzpyKzfOFpvg/SelCP802RItKYFPfSLRxJ3jf0bV19hvYOEktPAYW2nd7S8MBcyHYyacHKbISq5iTmDzG+ivnR+SZJv3USNTYVMz9rCQVSxd0cLlgsJauko79NnwOJbzDTyLooYolwz75gBOH2/xOMirpeEqRJrF/EQjWekJmGk9RtboXePu2rka+Xm51syBPhiXJAq0GfbfaFu9Nqs/e2Via/tEIM0lqvxfXQ6da6HrThsm5id4ClZFliOacRfrsPLRixS/IQYtksxghvJwbqOsbIsITail9Ayy4tKcogeEZiOO+4Ed264NSKmk7l3wKwJLAFjCFogBx8GE3OBz4pqcAn/ydA="
	time := "20160729142400"
	seq := "0001"
	secret := key + data + time + seq

	keys := []byte(key)
	h := hmac.New(sha256.New, keys)
	h.Write([]byte(secret))
	str := fmt.Sprintf("%x", hex.EncodeToString(h.Sum(nil)))
	fmt.Printf("%d-%s", len(str), str)
}

func TestMd5(t *testing.T) {
	str := fmt.Sprintf("%x", md5.Sum([]byte("0317ca0cebd4")))
	fmt.Println(len(str))
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
	origData = pkcs7Padding(origData, blockSize)
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
	origin, err := pkcs7Trimming(origin)
	if err != nil {
		return err.Error()
	}
	return string(origin)
}
