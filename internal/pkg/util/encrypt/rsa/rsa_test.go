package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"testing"
	"unicode/utf8"
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

func byHex(src string) []byte {
	decodeString, err := hex.DecodeString(src)
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return decodeString
}

func TestRsaCrypt(t *testing.T) {
	//bs, err := ioutil.ReadFile("2.prikey")
	//if err != nil {
	//	t.Fatalf("read private key file failed, err = %s", err.Error())
	//}
	//block, _ := pem.Decode(bs)
	//if block == nil {
	//	t.Fatalf("decode private key failed")
	//}
	//privt, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	////privt := pri.(*rsa.PrivateKey)
	//if err != nil {
	//	t.Fatalf("parse private key failed, err=%s", err)
	//}
	//testData := "hellow, afsaf, what the fuck ??????,世界还是美好的"
	//
	//testEncData, err := PrivateEncrypt(privt, []byte(testData))
	//if err != nil {
	//	t.Fatalf("private key encrypt failed, err=%s", err)
	//}
	//fmt.Printf("orig data=%s\nenc data=%s\n",
	//	testData, base64.StdEncoding.EncodeToString(testEncData))

	testData := "5321903be860e50a6568d24fe40572164a6d28b1a417b13120894c12683461e24c576e36715dfdc34f7f841052c4cd78395296500734b0024e4dbf20e4a3cc783d4d4da59d815b734184011797b47027771b4282ab1823e6c94aeb2ba6c1b2c292d2a95138818d6f90ba742029f94da9fdc5de8c1fd99687539f33b1e77f029e8d43d48cdbc77947f8f4eb16e9e3a970ae87a5516805a1aa844d07cfecbae92af4ba8e7972056f8215841e4a4233faaaa1c58e2c2e44e348b6c58b9943b2ded400c0527aaea11eef59206d9096bda1c05d2e649fe7460f0d0b7b7aaf9678a2a254438106ccfacdf2429799b6c7b510bff0c781ad8320ba960c58c2d2a8f16df1"
	testEncData := byHex(testData)

	bs, err := ioutil.ReadFile("2.pubkey")
	if err != nil {
		t.Fatalf("read certifacte file failed, err = %s", err.Error())
	}

	block, _ := pem.Decode(bs)
	if block == nil {
		t.Fatalf("decode certifacte key failed")
	}
	pu, err := x509.ParsePKIXPublicKey(block.Bytes)
	pub := pu.(*rsa.PublicKey)
	if err != nil {
		t.Fatalf("parse certifacte key failed, err=%s", err)
	}

	//cert, err := x509.ParseCertificate(block.Bytes)
	//if err != nil {
	//	t.Fatalf("parse certifacte key failed, err=%s", err)
	//}
	//pub := cert.PublicKey.(*rsa.PublicKey)
	testDecData, err := PublicDecrypt(pub, testEncData)
	if err != nil {
		t.Fatalf("public decrypt failed")
	}
	r, _ := utf8.DecodeRune(testDecData)
	fmt.Printf("dec data = %s\n", string(r))

	fmt.Printf("dec data = %s\n", string(testDecData))
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
