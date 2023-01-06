package community

import (
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"
	"testing"
	"time"
)

func TestParaseTime(t *testing.T) {
	stringTime := 1672897347
	x := time.Unix(int64(stringTime), 0).Format(timetool.TimeFormat)
	fmt.Println(x)
	//y, _ := time.ParseInLocation(timetool.TimeFormat, stringTime, time.Local)
	//fmt.Println(y)
}
