package event

import (
	"encoding/json"
	eevent "mio/internal/pkg/model/entity/event"
)

var DefaultEventTemplateService = EventTemplateService{}

type EventTemplateService struct {
}

func (srv EventTemplateService) ParseSetting(t eevent.EventTemplateType, setting eevent.EventTemplateSetting) (eevent.EventTemplateSettingInfo, error) {
	switch t {
	case eevent.EventTemplateTypeLCAER:
		lcaer := eevent.EventTemplateLCAERSetting{}
		err := json.Unmarshal([]byte(setting), &lcaer)
		if err != nil {
			return nil, err
		}
		return &lcaer, nil
	case eevent.EventTemplateTypeHC:
		hc := eevent.EventTemplateHCSetting{}
		err := json.Unmarshal([]byte(setting), &hc)
		if err != nil {
			return nil, err
		}
		return &hc, nil
	case eevent.EventTemplateTypeEEP:
		eep := eevent.EventTemplateEEPSetting{}
		err := json.Unmarshal([]byte(setting), &eep)
		if err != nil {
			return nil, err
		}
		return &eep, nil
	}
	return nil, nil
}
