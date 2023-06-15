package charge

import (
	"fmt"
	"math/rand"
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

//获取流水号订单号

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateSerialNumber() string {
	operatorID := "MA1G55M8X"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 18)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	result := string(b)
	return operatorID + result
}
