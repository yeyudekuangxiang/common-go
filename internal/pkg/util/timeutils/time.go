package timeutils

import "time"

func Now() Time {
	return Time{Time: time.Now()}
}

type Time struct {
	Time time.Time
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
func (t Time) Format() string {
	return t.Time.Format("2006-01-02 15:04:05")
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
