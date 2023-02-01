package admin

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	day1 := time.Now().Day()

	time1 := time.Date(2023, 1, 30, 0, 0, 0, 0, time.Local)
	day2 := time1.Day()
	time1.Unix()
	time1.Format("2006-01-02")
	fmt.Printf("day1 = %v\n", day1)
	fmt.Printf("day2 = %v\n", day2)

}
