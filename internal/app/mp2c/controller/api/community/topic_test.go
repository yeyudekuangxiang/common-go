package community

import (
	"bytes"
	"encoding/json"
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

func TestCheckSignupInfo(t *testing.T) {
	info := "[{\"type\":1,\"title\":\"姓名\",\"sort\":1},{\"type\":2,\"title\":\"性别\",\"sort\":2,\"options\":{\"option1\":\"男\",\"option2\":\"女\"}},{\"type\":4,\"title\":\"爱好\",\"sort\":3,\"options\":{\"option1\":\"唱\",\"option2\":\"跳\",\"option3\":\"rap\",\"option4\":\"篮球\"}},{\"type\":3,\"title\":\"备注\",\"sort\":4}]"
	infos := make([]interface{}, 0)
	err := json.Unmarshal([]byte(info), &infos)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%+v", infos)
}
