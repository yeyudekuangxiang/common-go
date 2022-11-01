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
	return Time{Time: t}
}

func Now() Time {
	return Time{Time: time.Now()}
}

type Time struct {
	time.Time
}

func (t *Time) StartOfDay() Time {
	return Time{Time: StartOfDay(t.Time)}
}
func (t *Time) EndOfDay() Time {
	return Time{Time: EndOfDay(t.Time)}
}
func (t *Time) StartOfWeek() Time {
	return Time{Time: StartOfWeek(t.Time)}
}
func (t *Time) EndOfWeek() Time {
	return Time{Time: EndOfWeek(t.Time)}
}
func (t *Time) StartOfMonth() Time {
	return Time{Time: StartOfMonth(t.Time)}
}
func (t *Time) EndOfMonth() Time {
	return Time{Time: EndOfMonth(t.Time)}
}
func (t *Time) Format(format string) string {
	return t.Time.Format(format)
}
func (t *Time) AddDay(day int) Time {
	return Time{Time: t.Time.AddDate(0, 0, day)}
}
func (t *Time) AddWeek(week int) Time {
	return Time{Time: t.Time.AddDate(0, 0, week*7)}
}
func (t *Time) AddMonth(month int) Time {
	return Time{Time: t.Time.AddDate(0, month, 0)}
}
func (t *Time) AddYear(year int) Time {
	return Time{Time: t.Time.AddDate(year, 0, 0)}
}
func (t *Time) String() string {
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
func (t *Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(fmt.Sprintf("\"\"")), nil
	}
	var stamp = fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(stamp), nil
}
func (t *Time) Value() (driver.Value, error) {
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
func (t *Time) Date() Date {
	return ToDate(t.Time)
}

func (t *Time) GetDiffDays(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	return int(t1.Sub(t2).Hours() / 24)
}

// 返回两日期相差 格式为 年 月 日
func (t *Time) SubTime(t1, t2 time.Time) string {
	y1 := t1.Year()
	y2 := t2.Year()
	m1 := int(t1.Month())
	m2 := int(t2.Month())
	d1 := t1.Day()
	d2 := t2.Day()

	//年差
	yearInterval := y1 - y2
	if m1 < m2 || m1 == m2 && d1 < d2 {
		yearInterval--
	}

	//月差
	monthInterval := (m1 + 12) - m2
	if d1 < d2 {
		monthInterval--
	}
	monthInterval %= 12
	
	//日差
	dayInterval := d1 - d2
	return fmt.Sprintf("%d年%d月%d日", yearInterval, monthInterval, dayInterval)
}

// 获取t1和t2的相差天数，单位：秒，0表同一天，正数表t1>t2，负数表t1<t2

/*func (t Time) GetDiffDaysBySecond(t1, t2 int64) int {
	time1 := time.Unix(t1, 0)
	time2 := time.Unix(t2, 0)
	// 调用上面的函数
	return Time.GetDiffDays(time1, time2)
}
*/
