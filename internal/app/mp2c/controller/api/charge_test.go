package api

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckTime(t *testing.T) {
	startTime, _ := time.Parse("2006-01-02", "2022-08-22")
	endTime, _ := time.Parse("2006-01-02", "2022-08-31")
	//fmt.Printf("start error : %s", err.Error())
	fmt.Println(time.Now())
	fmt.Println(startTime)
	fmt.Println(endTime)

	if time.Now().After(startTime) && time.Now().Before(endTime) {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}
