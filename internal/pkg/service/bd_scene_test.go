package service

import (
	"fmt"
	"mio/pkg/duiba/util"
	"testing"
)

func TestSign(t *testing.T) {
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
