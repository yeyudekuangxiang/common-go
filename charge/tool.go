package charge

import (
	"fmt"
	"time"
)

func TimeToDuration(startTimeStr string, endTimeStr string) string {
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		return ""
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		return ""
	}
	duration := endTime.Sub(startTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) - hours*60
	seconds := int(duration.Seconds()) - hours*3600 - minutes*60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
