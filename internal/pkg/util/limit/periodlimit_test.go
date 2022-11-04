package limit

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestNewPeriodLimit(t *testing.T) {
	now := time.Now()
	fmt.Println(now)
	name, offset := now.Zone()
	fmt.Println(name, offset)
	unix := now.Unix() + int64(offset)
	fmt.Println(unix)
	ti := 86400 - int(unix%int64(86400))
	fmt.Println(ti)
	duration, _ := time.ParseDuration(strconv.Itoa(ti) + "s")
	fmt.Println(duration.Hours())
}
