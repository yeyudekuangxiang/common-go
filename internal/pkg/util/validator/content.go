package validator

import (
	"bytes"
	"fmt"
	"github.com/medivhzhan/weapp/v3/security"
	"mio/internal/pkg/core/app"
)

func CheckMsgWithOpenId(openid, content string) error {
	length := len(content)
	if length > 2500 {
		s := []rune(content)
		var buffer bytes.Buffer
		for i, str := range s {
			buffer.WriteString(string(str))
			if i > 0 && (i+1)%2500 == 0 {
				params := &security.MsgSecCheckRequest{
					Content: buffer.String(),
					Version: 2,
					Scene:   2,
					Openid:  openid,
				}
				err := checkMsg(params)
				buffer.Reset()
				if err != nil {
					return err
				}
			}
		}
		params := &security.MsgSecCheckRequest{
			Content: buffer.String(),
			Version: 2,
			Scene:   2,
			Openid:  openid,
		}
		err := checkMsg(params)
		buffer.Reset()
		if err != nil {
			return err
		}
		return nil
	}
	//处理内容
	params := &security.MsgSecCheckRequest{
		Content: content,
		Version: 2,
		Scene:   2,
		Openid:  openid,
	}
	err := checkMsg(params)
	if err != nil {
		return err
	}
	return nil
}

func checkMsg(params *security.MsgSecCheckRequest) error {
	check, err := app.Weapp.NewSecurity().MsgSecCheck(params)
	if err != nil {
		return err
	}
	if check.ErrCode != 0 {
		return fmt.Errorf("check error: %s", check.ErrMSG)
	}
	if check.Result.Suggest != "pass" && check.Result.Label != 100 {
		return fmt.Errorf("check error: %s", "内容不合规，请重新输入")
	}
	return nil
}
