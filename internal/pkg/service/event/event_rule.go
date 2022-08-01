package event

import (
	eevent "mio/internal/pkg/model/entity/event"
	revent "mio/internal/pkg/repository/event"
	"strings"
)

var DefaultEventRuleService = EventRuleService{repo: revent.DefaultEventRuleRepository}

type EventRuleService struct {
	repo revent.EventRuleRepository
}

func (srv EventRuleService) GetFormattedEventRule(eventId string) (string, error) {
	list, err := srv.GetEventRuleList(GetEventRuleListParam{
		EventId: eventId,
	})
	if err != nil {
		return "", err
	}
	return srv.FormatEventRule(list), nil
}
func (srv EventRuleService) GetEventRuleList(param GetEventRuleListParam) ([]eevent.EventRule, error) {
	return srv.repo.GetEventRuleList(revent.GetEventRuleListBy{
		EventId: param.EventId,
	})
}
func (srv EventRuleService) FormatEventRule(list []eevent.EventRule) string {
	ht := strings.Builder{}
	for _, detail := range list {
		if strings.HasPrefix(detail.Content, "https://") || strings.HasPrefix(detail.Content, "http://") {
			ht.WriteString(`<p><img style="height:100%;width:100%" src="`)
			ht.WriteString(detail.Content)
			ht.WriteString(`"></img></p>`)
		} else {
			ht.WriteString(`<p>`)
			ht.WriteString(detail.Content)
			ht.WriteString(`</p>`)
		}
	}
	return ht.String()
}
func (srv EventRuleService) ParseEventContent(content string) []string {
	list := strings.Split(content, "\n")
	contents := make([]string, 0)
	strBuilder := strings.Builder{}
	for _, item := range list {
		if strings.HasPrefix(item, "https://") || strings.HasPrefix(item, "http://") {
			if strBuilder.Len() > 0 {
				contents = append(contents, strBuilder.String())
				strBuilder.Reset()
			}
			contents = append(contents, item)
		} else {
			if strBuilder.Len() > 0 {
				strBuilder.WriteString("\n")
			}
			strBuilder.WriteString(item)
		}
	}
	if strBuilder.Len() > 0 {
		contents = append(contents, strBuilder.String())
		strBuilder.Reset()
	}
	return contents
}
