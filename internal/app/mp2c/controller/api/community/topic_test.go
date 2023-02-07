package community

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCheckMsg(t *testing.T) {
	s := "你妈的哈哈哈你在搞笑吗"
	str := []rune(s)
	//res := ""
	//i := 0
	var buffer bytes.Buffer
	for i, content := range str {
		buffer.WriteString(string(content))
		//fmt.Printf("i%d r %s\n", i, buffer.String())
		if i > 0 && (i+1)%3 == 0 {
			fmt.Printf("=>(%d) '%v'\n", i, buffer.String())
			buffer.Reset()
		}
	}
	fmt.Printf("=>(0) '%v'\n", buffer.String())
}
