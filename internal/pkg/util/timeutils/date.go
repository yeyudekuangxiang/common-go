package timeutils

import (
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const DateFormat = "2006-01-02"
const FullDateFormat = "2006-01-02 15:04:05"

func ToDate(t time.Time) Date {
	return Date{time: StartOfDay(t)}
}
func NowDate() Date {
	return Date{time: StartOfDay(time.Now())}
}

type Date struct {
	time time.Time
}

func (d Date) Time() time.Time {
	return d.time
}
func (d Date) Format(format string) string {
	return d.time.Format(format)
}
func (d Date) String() string {
	return d.time.Format("2006-01-02")
}
func (d Date) FullString() string {
	return d.time.Format(FullDateFormat)
}
func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		return nil
	}

	ti, err := time.Parse(DateFormat, strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	d.time = ti
	return nil
}
func (d Date) MarshalJSON() ([]byte, error) {
	if d.time.IsZero() {
		return []byte("\"\""), nil
	}
	var stamp = fmt.Sprintf("\"%s\"", d.time.Format(DateFormat))
	return []byte(stamp), nil
}
func (d Date) Value() (driver.Value, error) {
	if d.time.IsZero() {
		return nil, nil
	}
	return d.time, nil
}
func (d *Date) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return errors.New("Date type error")
	}
	d.time = t
	return nil
}
func (d Date) StartOfWeek() Date {
	return Date{time: StartOfWeek(d.time)}
}
func (d Date) EndOfWeek() Date {
	return Date{time: EndOfWeek(d.time)}
}
func (d Date) StartOfMonth() Date {
	return Date{time: StartOfMonth(d.time)}
}
func (d Date) EndOfMonth() Date {
	return Date{time: EndOfMonth(d.time)}
}
func (d Date) AddDay(day int) Date {
	return Date{time: d.time.AddDate(0, 0, day)}
}
func (d Date) AddWeek(week int) Date {
	return Date{time: d.time.AddDate(0, 0, week*7)}
}
func (d Date) AddMonth(month int) Date {
	return Date{time: d.time.AddDate(0, month, 0)}
}
func (d Date) AddYear(year int) Date {
	return Date{time: d.time.AddDate(year, 0, 0)}
}
