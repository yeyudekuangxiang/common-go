package encrypttool

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

type getToken struct {
	OperatorSecret string `json:"OperatorSecret,omitempty"`
	OperatorID     string `json:"OperatorID,omitempty"`
	Sig            string `json:"Sig,omitempty"`
	Data           string `json:"Data,omitempty"`
	TimeStamp      string `json:"TimeStamp,omitempty"`
	Seq            string `json:"Seq,omitempty"`
}

func TestHMAC_Md5(t *testing.T) {
	key := "1234567890abcdef"
	data := "l7B0BSEjFdzpyKzfOFpvg/SelCP802RItKYFPfSLRxJ3jf0bV19hvYOEktPAYW2nd7S8MBcyHYyacHKbISq5iTmDzG+ivnR+SZJv3USNTYVMz9rCQVSxd0cLlgsJauko79NnwOJbzDTyLooYolwz75gBOH2/xOMirpeEqRJrF/EQjWekJmGk9RtboXePu2rka+Xm51syBPhiXJAq0GfbfaFu9Nqs/e2Via/tEIM0lqvxfXQ6da6HrThsm5id4ClZFliOacRfrsPLRixS/IQYtksxghvJwbqOsbIsITail9Ayy4tKcogeEZiOO+4Ed264NSKmk7l3wKwJLAFjCFogBx8GE3OBz4pqcAn/ydA="
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(data))
	str := hex.EncodeToString(h.Sum([]byte("")))
	str2 := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Printf("%d-%s\n", len(str), str)
	fmt.Printf("%d-%s\n", len(str2), str2)
	fmt.Printf("%d-%s\n", len("745166E8C43C84D37FFECOF529C4136F"), "745166E8C43C84D37FFECOF529C4136F")

}

func TestMd5(t *testing.T) {
	Md5 := md5.New()
	Md5.Write([]byte("123"))
	md5Data := Md5.Sum([]byte(""))
	str := hex.EncodeToString(md5Data)
	fmt.Printf("%s\n", str)
}

func TestAES(t *testing.T) {
	decryptCode := aesDecrypt("j5tJ74cKFiGJ65Ot7NaSyZQoaYNUpSYy7hVWul9Yw26tXyLZb7F2Vf+58kGMk6GUfUzR6WVJn7asnFnL7UfoNg==", "a2164ada0026ccf7", "82c91325e74bef0f")
	fmt.Println("解密结果：", decryptCode)
	data := getToken{
		OperatorSecret: "acb93539fc9bg78k",
		OperatorID:     "MA1G55M81",
	}
	marshal, _ := json.Marshal(data)
	fmt.Println("json marshal：", string(marshal))
	decryptCode1 := aesEncrypt(string(marshal), "a2164ada0026ccf7", "82c91325e74bef0f")
	fmt.Println("加密结果：", decryptCode1)

	decryptCode2 := aesEncrypt(decryptCode, "a2164ada0026ccf7", "82c91325e74bef0f")
	fmt.Println("加密结果：", decryptCode2)
}

func aesEncrypt(orig string, key, iv string) string {
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

func aesDecrypt(cryted string, key, iv string) string {
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
		return err.Error()
	}
	return string(origin)
}
