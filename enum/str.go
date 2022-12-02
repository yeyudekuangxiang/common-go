package enum

import (
	"database/sql/driver"
	"encoding/json"
)

var _ IEnumStrStatus = (*StrStatus)(nil)

type IEnumStrStatus interface {
	IDOAOIIJQWQOWJEOIHIOHSDIDHOSD()
	Text() string
	RealText() string
	Status() string
	Others() []string
	driver.Valuer
	json.Marshaler
}
type StrStatus struct {
	status   string
	text     string
	realText string
	others   []string
}

func (e StrStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.status + `"`), nil
}

func (e StrStatus) Value() (driver.Value, error) {
	return e.status, nil
}
func (e StrStatus) RealText() string {
	return e.realText
}
func (e StrStatus) Status() string {
	return e.status
}
func (e StrStatus) Text() string {
	return e.text
}
func (e StrStatus) Others() []string {
	return e.others
}
func (e StrStatus) IDOAOIIJQWQOWJEOIHIOHSDIDHOSD() {}
func NewEnumStrStatus(status string, text string, realText string, others ...string) IEnumStrStatus {
	return StrStatus{status: status, text: text, realText: realText, others: others}
}

type EnumStrStatusList []IEnumStrStatus

func (l EnumStrStatusList) Map() map[string]IEnumStrStatus {
	m := make(map[string]IEnumStrStatus)
	for _, item := range l {
		m[item.Status()] = item
	}
	return m
}
func (l EnumStrStatusList) Exist(status string) bool {
	_, ok := l.Map()[status]
	return ok
}
func (l EnumStrStatusList) List() []IEnumStrStatus {
	return l
}
func (l EnumStrStatusList) Find(status string) (enumStatus IEnumStrStatus, exist bool) {
	enumStatus, exist = l.Map()[status]
	return
}
