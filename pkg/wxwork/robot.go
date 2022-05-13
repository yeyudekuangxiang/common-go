package wxwork

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func SendRobotMessage(key string, v interface{}) error {
	var msgType MsgType
	switch v.(type) {
	case Text:
		msgType = MsgTypeText
	case Markdown:
		msgType = MsgTypeMarkdown
	case Image:
		msgType = MsgTypeImage
	case News:
		msgType = MsgTypeNews
	case File:
		msgType = MsgTypeFile
	default:
		return errors.New("unknown msgtype")
	}

	data, err := json.Marshal(map[string]interface{}{
		"msgtype":       msgType,
		string(msgType): v,
	})
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	resp, err := http.Post("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="+key, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))
	return nil
}
