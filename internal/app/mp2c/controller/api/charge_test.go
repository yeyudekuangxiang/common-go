package api

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckTime(t *testing.T) {
	startTime, err := time.Parse("2006-01-02", "2022-08-22")
	if err != nil {
		fmt.Printf("err : %s\n", err.Error())
	}
	endTime, err := time.Parse("2006-01-02", "2022-08-31")
	if err != nil {
		fmt.Printf("err : %s\n", err.Error())
	}
	fmt.Println(time.Now().String())
	fmt.Println(startTime.String())
	fmt.Println(endTime.String())

	fmt.Println(time.Now().After(startTime))
	fmt.Println(time.Now().Before(endTime))
}
