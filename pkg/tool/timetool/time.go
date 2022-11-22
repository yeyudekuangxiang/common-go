package timetool

import (
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

func ToTime(t time.Time) Time {
	return Time{Time: t}
}

func Now() Time {
	return Time{Time: time.Now()}
}

type Time struct {
	time.Time
}

func (t Time) StartOfDay() Time {
	return Time{Time: StartOfDay(t.Time)}
}
func (t Time) EndOfDay() Time {
	return Time{Time: EndOfDay(t.Time)}
}
func (t Time) StartOfWeek() Time {
	return Time{Time: StartOfWeek(t.Time)}
}
func (t Time) EndOfWeek() Time {
	return Time{Time: EndOfWeek(t.Time)}
}
func (t Time) StartOfMonth() Time {
	return Time{Time: StartOfMonth(t.Time)}
}
func (t Time) EndOfMonth() Time {
	return Time{Time: EndOfMonth(t.Time)}
}
func (t Time) Format(format string) string {
	if t.IsZero() {
		return ""
	}
	return t.Time.Format(format)
}
func (t Time) AddDay(day int) Time {
	return Time{Time: t.Time.AddDate(0, 0, day)}
}
func (t Time) AddWeek(week int) Time {
	return Time{Time: t.Time.AddDate(0, 0, week*7)}
}
func (t Time) AddMonth(month int) Time {
	return Time{Time: t.Time.AddDate(0, month, 0)}
}
func (t Time) AddYear(year int) Time {
	return Time{Time: t.Time.AddDate(year, 0, 0)}
}
func (t Time) Add(duration time.Duration) Time {
	return Time{Time: t.Time.Add(duration)}
}
func (t Time) String() string {
	return t.Time.Format(TimeFormat)
}
func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		return nil
	}
	ti, err := time.Parse(TimeFormat, strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	t.Time = ti
	return nil
}
func (t Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(fmt.Sprintf("\"\"")), nil
	}
	var stamp = fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(stamp), nil
}
func (t Time) Value() (driver.Value, error) {
	if t.Time.IsZero() {
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
func (t Time) Date() Date {
	return ToDate(t.Time)
}
func Parse(layout string, val string) (Time, error) {
	if val == "" {
		return Time{}, nil
	}
	t, err := time.Parse(layout, val)
	return Time{
		Time: t,
	}, err
}
func MustParse(layout string, val string) Time {
	if val == "" {
		return Time{}
	}
	t, err := time.Parse(layout, val)
	if err != nil {
		panic(err)
	}
	return Time{
		Time: t,
	}
}
func (t Time) IfNotZero(f func(t Time) Time) Time {
	return f(t)
}
func UnixMilli(msec int64) Time {
	return UnixMilliZero(msec, true)
}

func UnixMilliZero(msec int64, zeroTime bool) Time {
	if msec == 0 && zeroTime {
		return Time{}
	}
	return Time{Time: time.UnixMilli(msec)}
}
