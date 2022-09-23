package platform

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	//createTime, _ := time.Parse(time.RFC3339Nano, "2022-08-30T14:33:49.807482+08:00")
	//updateTime, _ := time.Parse(time.RFC3339Nano, "2022-09-02T13:16:56.369145+08:00")
	//fmt.Printf("%s\n", createTime)
	//fmt.Printf("%s\n", updateTime)
	//fmt.Printf("%v\n", updateTime.After(createTime))
	startTime, _ := time.Parse("2006-01-02", "2022-09-23")
	endTime, _ := time.Parse("2006-01-02", "2022-10-01")
	fmt.Printf("%v\n", startTime.Format("2006-01-02"))
	fmt.Printf("%s\n", endTime)
	fmt.Printf("%s\n", startTime)
	fmt.Printf("%s\n", endTime.String())
	fmt.Printf("%v\n", time.Now())

}
