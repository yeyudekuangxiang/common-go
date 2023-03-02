package chuanglan

import (
	"testing"
)

func TestSign(t *testing.T) {
	NewSmsClient("1", "2").SendV2("1", "w")
}
