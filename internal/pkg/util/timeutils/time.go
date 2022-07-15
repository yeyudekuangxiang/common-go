package timeutils

import (
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

func ToTime(t time.Time) Time {
	return Time{time: t}
}
func NowTime() Time {
	return Time{time: time.Now()}
}

type Time struct {
	time time.Time
}

func (t Time) Time() time.Time {
	return t.time
}
func (t Time) StartOfDay() Time {
	return Time{time: StartOfDay(t.time)}
}
func (t Time) EndOfDay() Time {
	return Time{time: EndOfDay(t.time)}
}
func (t Time) StartOfWeek() Time {
	return Time{time: StartOfWeek(t.time)}
}
func (t Time) EndOfWeek() Time {
	return Time{time: EndOfWeek(t.time)}
}
func (t Time) StartOfMonth() Time {
	return Time{time: StartOfMonth(t.time)}
}
func (t Time) EndOfMonth() Time {
	return Time{time: EndOfMonth(t.time)}
}
func (t Time) Format(format string) string {
	return t.time.Format(format)
}
func (t Time) AddDay(day int) Time {
	return Time{time: t.time.AddDate(0, 0, day)}
}
func (t Time) AddWeek(week int) Time {
	return Time{time: t.time.AddDate(0, 0, week*7)}
}
func (t Time) AddMonth(month int) Time {
	return Time{time: t.time.AddDate(0, month, 0)}
}
func (t Time) AddYear(year int) Time {
	return Time{time: t.time.AddDate(year, 0, 0)}
}
func (t Time) String() string {
	return t.time.Format(TimeFormat)
}
func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		return nil
	}
	ti, err := time.Parse(TimeFormat, strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	t.time = ti
	return nil
}
func (t Time) MarshalJSON() ([]byte, error) {
	if t.time.IsZero() {
		return []byte(fmt.Sprintf("\"\"")), nil
	}
	var stamp = fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(stamp), nil
}
func (t Time) Value() (driver.Value, error) {
	if t.time.IsZero() {
		return nil, nil
	}
	return t.time, nil
}
func (t *Time) Scan(value interface{}) error {
	ti, ok := value.(time.Time)
	if !ok {
		return errors.New("Time type error")
	}
	t.time = ti
	return nil
}
func (t Time) Date() Date {
	return ToDate(t.time)
}
