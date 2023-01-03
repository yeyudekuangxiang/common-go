package timetool

import (
	"time"
)

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
func EndOfDay(t time.Time) time.Time {
	t = t.AddDate(0, 0, 1)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, -1, t.Location())
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
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-1)
}
