package jhx

import (
	"fmt"
	"testing"
	"time"
)

func TestExpireTime(t *testing.T) {
	expireTime, _ := time.Parse("2006-01-02", "2022-12-12")
	expireTimeUnix := expireTime.UnixMilli()
	fmt.Println(expireTimeUnix)
}
