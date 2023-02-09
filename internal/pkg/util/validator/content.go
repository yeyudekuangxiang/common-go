package validator

import (
	"bytes"
	"errors"
	"github.com/medivhzhan/weapp/v3/security"
	"mio/internal/pkg/core/app"
	"strings"
)

var labelMap = map[int]string{
	10001: "广告",
	20001: "时政",
	20002: "色情",
	20003: "辱骂",
	20006: "违法犯罪",
	20008: "欺诈",
	20012: "低俗",
	20013: "版权",
	21000: "其他",
}

func CheckMsgWithOpenId(openid, content string) error {
	content = strings.ReplaceAll(content, " ", "")
	content = strings.ReplaceAll(content, "\n", "")
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
	var check *security.MsgSecCheckResponse
	err := app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		check, err = app.Weapp.NewSecurity().MsgSecCheck(params)
		if err != nil {
			return false, err
		}
		return app.Weapp.IsExpireAccessToken(check.ErrCode)
	}, 1)

	if err != nil {
		return err
	}

	if check.ErrCode != 0 {
		return errors.New(check.ErrMSG)
	}

	if check.Result.Suggest != "pass" && check.Result.Label != 100 {
		return errors.New("内容不合规: " + labelMap[check.Result.Label])
	}

	return nil
}

func CheckMediaWithOpenId(openid, mediaUrl string) error {
	req := &security.MediaCheckAsyncRequest{
		MediaUrl:  mediaUrl + "?x-oss-process=image/resize,m_fixed,h_100,w_100",
		MediaType: 2,
		Version:   2,
		Openid:    openid,
		Scene:     3,
	}
	var resp *security.MediaCheckAsyncResponse
	err := app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		resp, err = app.Weapp.NewSecurity().MediaCheckAsync(req)
		if err != nil {
			return false, err
		}
		return app.Weapp.IsExpireAccessToken(resp.ErrCode)
	}, 1)

	if err != nil {
		return err
	}

	if resp.ErrCode != 0 {
		return errors.New(resp.ErrMSG)
	}

	return nil
}
