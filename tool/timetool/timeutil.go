package timetool

import (
	"database/sql"
	"time"
)

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}
func EndOfDay(t time.Time) time.Time {
	t = t.AddDate(0, 0, 1)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, -1, time.Local)
}
func StartOfWeek(t time.Time) time.Time {
	day := t.Weekday() - 1
	if day == -1 {
		day = 6
	}
	return StartOfDay(t.Add(time.Hour * 24 * time.Duration(-day)))
}
func EndOfWeek(t time.Time) time.Time {
	day := (7 - t.Weekday()) % 7
	return EndOfDay(t.Add(time.Hour * 24 * time.Duration(day)))
}
func StartOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
}
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-1)
}
func Format(t time.Time, format string, zeroStr string) string {
	if t.IsZero() {
		return zeroStr
	}

	return t.Format(format)
}
func UnixMilliNullTIme(msec int64) sql.NullTime {
	if msec <= 0 {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Valid: true,
		Time:  time.UnixMilli(msec),
	}
}
