package event

import (
	"encoding/json"
	"github.com/pkg/errors"
	"mio/internal/pkg/util"
	"reflect"
	"strings"
)

// EventTemplateSetting 公益活动证书模版配置信息 json字符串
type EventTemplateSetting string

// EventTemplateType 公益活动证书模版类型
type EventTemplateType string

const (
	EventTemplateTypeEEP   EventTemplateType = "EEP"   //生态环保 Ecological environmental protection
	EventTemplateTypeHC    EventTemplateType = "HC"    //人文关怀 humanistic concern
	EventTemplateTypeLCAER EventTemplateType = "LCAER" //低碳减排 Low carbon and emission reduction
)

// Text 公益活动证书模版类型对应的中文
func (et EventTemplateType) Text() string {
	switch et {
	case EventTemplateTypeEEP:
		return "生态环保"
	case EventTemplateTypeHC:
		return "人文关怀"
	case EventTemplateTypeLCAER:
		return "低碳减排"
	}
	return ""
}

var (
	// EventTenmplateList 公益活动证书模版类型列表
	EventTenmplateList = []EventTemplateType{
		EventTemplateTypeEEP,
		EventTemplateTypeHC,
		EventTemplateTypeLCAER,
	}
)

type EventTemplateSettingInfo interface {
	Setting() EventTemplateSetting
}

type EventTemplateEEPSetting struct {
	Recipient    string `json:"recipient"  setting:"label:被捐助机构;desc:被捐助机构名称"`
	Money        string `json:"money" setting:"label:贡献金额;desc:贡献金额"`
	Desc         string `json:"desc" setting:"label:公益描述;desc:公益描述"`
	Image        string `json:"image" setting:"label:公益图片;desc:公益图片;type:url"`
	Organization string `json:"organization" setting:"label:项目方;desc:项目方"`
}

func (e EventTemplateEEPSetting) Setting() EventTemplateSetting {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return EventTemplateSetting(data)
}

type EventTemplateHCSetting struct {
	Project      string `json:"project" setting:"label:公益项目名称;desc:公益项目名称"`
	Money        string `json:"money" setting:"label:贡献金额;desc:贡献金额"`
	Desc         string `json:"desc" setting:"label:公益描述;desc:公益描述"`
	Image        string `json:"image" setting:"label:公益图片;desc:公益图片;type:url"`
	Organization string `json:"organization" setting:"label:项目方;"`
}

func (e EventTemplateHCSetting) Setting() EventTemplateSetting {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return EventTemplateSetting(data)
}

type EventTemplateLCAERSetting struct {
	Project string `json:"project" setting:"label:公益项目名称;desc:公益项目名称"`
	Money   string `json:"money" setting:"label:贡献金额;desc:贡献金额"`
	Desc    string `json:"desc" setting:"label:公益描述;desc:公益描述"`
	Image   string `json:"image" setting:"label:公益图片;desc:公益图片;type:url"`
}

func (e EventTemplateLCAERSetting) Setting() EventTemplateSetting {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return EventTemplateSetting(data)
}

type Field struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Label string      `json:"label"`
	Desc  string      `json:"desc"`
	Type  string      `json:"type"`
}

func Parse(v interface{}) ([]Field, error) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("不支持的类型")
	}
	val := reflect.ValueOf(v)
	list := make([]Field, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(0)
		tag := field.Tag.Get("setting")
		tags := strings.Split(tag, ";")
		tagMap := make(map[string]string)
		for _, tag := range tags {
			tagKV := strings.Split(tag, ":")
			if len(tagKV) == 2 {
				tagMap[tagKV[0]] = tagKV[1]
			}
		}
		list = append(list, Field{
			Key:   field.Name,
			Value: val.Field(i).Interface(),
			Label: util.Ternary(tagMap["label"] != "", tagMap["label"], field.Name).String(),
			Desc:  util.Ternary(tagMap["desc"] != "", tagMap["desc"], field.Name).String(),
			Type:  util.Ternary(tagMap["type"] != "", tagMap["type"], "string").String(),
		})
	}
	return list, nil
}
