package enum

import "database/sql/driver"

var _ IEnumIntStatus = (*IntStatus)(nil)

type IEnumIntStatus interface {
	IDOAOIIJQWQOWJEOIHIOHSDIDHOSD()
	Text() string
	RealText() string
	Status() int64
	Others() []string
	driver.Valuer
}

type IntStatus struct {
	status   int64
	text     string
	realText string
	others   []string
}

func (e IntStatus) Value() (driver.Value, error) {
	return e.status, nil
}
func (e IntStatus) RealText() string {
	return e.realText
}
func (e IntStatus) Status() int64 {
	return e.status
}
func (e IntStatus) Text() string {
	return e.text
}
func (e IntStatus) Others() []string {
	return e.others
}
func (e IntStatus) IDOAOIIJQWQOWJEOIHIOHSDIDHOSD() {}
func NewEnumIntStatus(status int64, text string, realText string, others ...string) IEnumIntStatus {
	return IntStatus{status: status, text: text, realText: realText, others: others}
}

type EnumIntStatusList []IEnumIntStatus

func (l EnumIntStatusList) Map() map[int64]IEnumIntStatus {
	m := make(map[int64]IEnumIntStatus)
	for _, item := range l {
		m[item.Status()] = item
	}
	return m
}
func (l EnumIntStatusList) Exist(status int64) bool {
	_, ok := l.Map()[status]
	return ok
}
func (l EnumIntStatusList) List() []IEnumIntStatus {
	return l
}
func (l EnumIntStatusList) Find(status int64) IEnumIntStatus {
	return l.Map()[status]
}
