package timeutils

import "time"

const defaultTimeFormat = "2006-01-02 15:04:05"
const defaultDateFormat = "2006-01-02"

func Now() Time {
	return Time{
		Time:   time.Now(),
		Format: defaultTimeFormat,
	}
}

type Time struct {
	Time   time.Time
	Format string
}

func (t Time) StartOfDay() Time {
	return Time{Time: StartOfDay(t.Time), Format: t.Format}
}
func (t Time) EndOfDay() Time {
	return Time{Time: EndOfDay(t.Time), Format: t.Format}
}
func (t Time) String() string {
	return t.Time.Format(t.Format)
}
