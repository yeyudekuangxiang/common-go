package service

import (
	"fmt"
	util2 "mio/internal/pkg/util/encrypt"
	"mio/pkg/duiba/util"
	"testing"
)

func TestCheckSign(t *testing.T) {
	sign := "3719139c1ea8a5e614b4927eb6f0a21a"
	key := "fastElectricity" + "#" + "18301939833" + "#" + "MA005DBW1220819154416041110" + "#" + "1.00" + "#" + "0qscxr0cebd4"
	encryptK := util.Md5(key)
	fmt.Println(encryptK)
	fmt.Println(sign)

	if encryptK == sign {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}

func TestSign(t *testing.T) {
	str := "lvmiao#13083605153#MA005DBW1220819154416041110#1.00#0317ca0cebd4"
	fmt.Println("localSignStr", str)
	localSign := util2.Md5(str)
	fmt.Printf("Sign: %s", localSign)
}
