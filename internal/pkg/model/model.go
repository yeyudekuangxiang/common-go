package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"
const dateFormat = "2006-01-02"

type Model struct {
	Id int64 `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
}

func NewTime() Time {
	return Time{
		time.Now(),
	}
}

//json时格式化日期为2016-01-02 15:04:05的格式
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		return nil
	}
	ti, err := time.Parse(timeFormat, strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	t.Time = ti
	return nil
}
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(fmt.Sprintf("\"\"")), nil
	}
	var stamp = fmt.Sprintf("\"%s\"", t.Format(timeFormat))
	return []byte(stamp), nil
}
func (t Time) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}
func (t *Time) Scan(value interface{}) error {
	ti, ok := value.(time.Time)
	if !ok {
		return errors.New("Time type error")
	}
	t.Time = ti
	return nil
}
func (t *Time) Date() Date {
	return Date{Time: t.Time}
}
func (t *Time) String() string {
	return t.Format(timeFormat)
}
func (t Time) StartOfDay() Time {
	return Time{Time: time.Date(t.Time.Year(), t.Time.Month(), t.Time.Day(), 0, 0, 0, 0, time.Local)}
}
func (t Time) EndOfDay() Time {
	t2 := Time{Time: t.Time.Add(time.Hour * 24)}
	return Time{Time: t2.StartOfDay().Add(-time.Nanosecond)}
}

func NewDate(date string) (Date, error) {
	t, err := time.Parse(dateFormat, date)
	if err != nil {
		return Date{}, err
	}
	return Date{
		Time: t,
	}, nil
}

//json时格式化时间为2016-01-02的格式
type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		return nil
	}

	ti, err := time.Parse(dateFormat, strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	d.Time = ti
	return nil
}
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("\"\""), nil
	}
	var stamp = fmt.Sprintf("\"%s\"", d.Time.Format(dateFormat))
	return []byte(stamp), nil
}
func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}
func (d *Date) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return errors.New("Date type error")
	}
	d.Time = t
	return nil
}
func (d Date) String() string {
	return d.Time.Format(dateFormat)
}
func (d Date) FullString() string {
	return d.Time.Format(dateFormat) + " 00:00:00"
}

//json时格式化日期为2016-01-02 15:04:05的格式 并且获取值时如果为初始时间则赋值当前时间
type CreatedTime struct {
	time.Time
}

func (ct *CreatedTime) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		return nil
	}
	ti, err := time.Parse(timeFormat, strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	ct.Time = ti
	return nil
}
func (ct CreatedTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", ct.Format(timeFormat))
	return []byte(stamp), nil
}
func (ct CreatedTime) Value() (driver.Value, error) {
	if ct.IsZero() {
		return time.Now().Format(timeFormat), nil
	}
	return ct.Format(timeFormat), nil
}
func (ct *CreatedTime) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return errors.New("CreatedTime type error")
	}
	ct.Time = t
	return nil
}

//json时格式化日期为2016-01-02 15:04:05的格式 并且获取值时如果为初始时间则赋值当前时间
type UpdatedTime struct {
	time.Time
}

func (ut *UpdatedTime) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		return nil
	}
	ti, err := time.Parse(timeFormat, strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	ut.Time = ti
	return nil
}
func (ut UpdatedTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", ut.Format(timeFormat))
	return []byte(stamp), nil
}
func (ut UpdatedTime) Value() (driver.Value, error) {
	return time.Now().Format(timeFormat), nil
}
func (ut *UpdatedTime) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return errors.New("CreatedTime type error")
	}
	ut.Time = t
	return nil
}

//存入数据库时转换成字符串形式,以英文逗号隔开
type ArrayString []string

func (as ArrayString) MarshalJSON() ([]byte, error) {
	if len(as) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(as))
}
func (as ArrayString) Value() (driver.Value, error) {
	if as == nil {
		return "", nil
	}
	return strings.Join(as, ","), nil
}
func (as *ArrayString) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return errors.New("ArrayString type error")
	}
	if len(v) > 0 {
		*as = strings.Split(v, ",")
	} else {
		*as = make([]string, 0)
	}
	return nil
}

type StrToInt int

func (s *StrToInt) UnmarshalJSON(data []byte) error {
	d, err := strconv.Atoi(strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	*s = StrToInt(d)
	return nil
}
func (s StrToInt) Value() (driver.Value, error) {
	return int(s), nil
}
func (s *StrToInt) Scan(value interface{}) error {
	t, ok := value.(int)
	if !ok {
		return errors.New("StrToInt type error")
	}
	*s = StrToInt(t)
	return nil
}
func (s StrToInt) Int() int {
	return int(s)
}

// 用户点赞
type UserLike struct {
	Model
	UserId     int64  `gorm:"not null;uniqueIndex:idx_user_like_unique;" json:"userId" form:"userId"`                                            // 用户
	EntityType string `gorm:"not null;size:32;uniqueIndex:idx_user_like_unique;index:idx_user_like_entity;" json:"entityType" form:"entityType"` // 实体类型
	EntityId   int64  `gorm:"not null;uniqueIndex:idx_user_like_unique;index:idx_user_like_entity;" json:"topicId" form:"topicId"`               // 实体编号
	CreateTime int64  `json:"createTime" form:"createTime"`                                                                                      // 创建时间
}

type NullString string

func (d NullString) Value() (driver.Value, error) {
	if d == "" {
		return nil, nil
	}
	return string(d), nil
}
func (d *NullString) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	t, ok := value.(string)
	if !ok {
		return errors.New("NullString type error")
	}
	*d = NullString(t)
	return nil
}

type NullInt int64

func (d NullInt) Value() (driver.Value, error) {
	if d == 0 {
		return nil, nil
	}
	return int64(d), nil
}
func (d *NullInt) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	t, ok := value.(int64)
	if !ok {
		return errors.New("NullInt type error")
	}
	*d = NullInt(t)
	return nil
}
