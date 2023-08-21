package gaode

import (
	"testing"
)

func TestName(t *testing.T) {
	ip, err := NewMapClient("456ddc091c99e74da11162e654a19983").LocationIp("36.98.227.232")
	if err != nil {
		return
	}

	println(ip)

}
