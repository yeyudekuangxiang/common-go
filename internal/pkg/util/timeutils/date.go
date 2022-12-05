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
	return Date{Time: StartOfDay(t)}
}
func NowDate() Date {
	return Date{Time: StartOfDay(time.Now())}
}

type Date struct {
	time.Time
}

func (d Date) Format(format string) string {
	return d.Time.Format(format)
}
func (d Date) String() string {
	return d.Time.Format("2006-01-02")
}
func (d Date) FullString() string {
	return d.Time.Format(FullDateFormat)
}
func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		return nil
	}

	ti, err := time.Parse(DateFormat, strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	d.Time = ti
	return nil
}
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("\"\""), nil
	}
	var stamp = fmt.Sprintf("\"%s\"", d.Time.Format(DateFormat))
	return []byte(stamp), nil
}
func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
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
func (d Date) StartOfWeek() Date {
	return Date{Time: StartOfWeek(d.Time)}
}
func (d Date) EndOfWeek() Date {
	return Date{Time: EndOfWeek(d.Time)}
}
func (d Date) StartOfMonth() Date {
	return Date{Time: StartOfMonth(d.Time)}
}
func (d Date) EndOfMonth() Date {
	return Date{Time: EndOfMonth(d.Time)}
}
func (d Date) AddDay(day int) Date {
	return Date{Time: d.Time.AddDate(0, 0, day)}
}
func (d Date) AddWeek(week int) Date {
	return Date{Time: d.Time.AddDate(0, 0, week*7)}
}
func (d Date) AddMonth(month int) Date {
	return Date{Time: d.Time.AddDate(0, month, 0)}
}
func (d Date) AddYear(year int) Date {
	return Date{Time: d.Time.AddDate(year, 0, 0)}
}