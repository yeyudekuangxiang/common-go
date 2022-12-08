package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"testing"
)

func assertArrayEqual(t *testing.T, a []byte, b []byte) {
	fmt.Printf("array1: %v, araray2: %v\n", a, b)
	if len(a) != len(b) {
		t.Fatal("array not equal!")
		return
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Fatalf("array not equal at position %d, %d != %d", i, a[i], b[i])
			return
		}
	}
}

func hexDecode(src string) []byte {
	decodeString, err := hex.DecodeString(src)
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return decodeString
}

func TestRsaCrypt(t *testing.T) {
	sign := "652d45759c8a42183bc52be15aad4998b438b40abc6f2a795d3f5d2e9e047f7cc8f4dbeb203de9366ddd6d7be60e04f8d7eced010662eba024a781a3ccb0a102e436fa56b5f32ca18e67933ceb4c29a097394a5139164a28249f6b7005b07ed265ccf8495abf13934a07f8c17e46216d41e962105bef3805c81adf4b35910bb2968d577f13b72a6a89ecb9fee5ad722a80a2bee5c69c89eb33dfd3378ec5caa84ebbc33d9caba2f94a954aa43b71bde0cd4e1071df9cc613fc9291f195572f0dd90e4caf803e33db9d7334d40d887ddda3159674b9d41c3f3c1d90b7377acb506f841ff01887fd4c30ee534fe8deb04e8a39b47d344ab51c6bfd6b743d601574"
	hexByte := hexDecode(sign)
	fmt.Printf("hexByte:%v\n", hexByte)
	//公钥
	bs, err := ioutil.ReadFile("2.key")
	if err != nil {
		fmt.Printf("read certifacte file failed, err = %s", err.Error())
	}
	block, _ := pem.Decode(bs)
	if block == nil {
		fmt.Printf("decode certifacte key failed")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}

	decrypt, err := PublicDecrypt(pub.(*rsa.PublicKey), hexByte)
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}

	hexString := hex.EncodeToString(decrypt)
	fmt.Printf("hexString:%v\n", hexString)
	fmt.Printf("decryptString: %s\n", string(decrypt))
}

func TestHex(t *testing.T) {
	// 转换的用的 byte数据
	byte_data := []byte(`测试数据`)
	// 将 byte 装换为 16进制的字符串
	hex_string_data := hex.EncodeToString(byte_data)
	// byte 转 16进制 的结果
	fmt.Println(hex_string_data)

	/* ====== 分割线 ====== */

	// 将 16进制的字符串 转换 byte
	hex_data, _ := hex.DecodeString(hex_string_data)
	// 将 byte 转换 为字符串 输出结果
	fmt.Println(string(hex_data))
}
