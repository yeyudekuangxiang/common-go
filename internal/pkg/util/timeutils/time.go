package timeutils

import "time"

type Time struct {
	Time time.Time
}

func (t Time) StartOfDay() Time {
	return Time{Time: StartOfDay(t.Time)}
}
func (t Time) EndOfDay() Time {
	return Time{Time: EndOfDay(t.Time)}
}
