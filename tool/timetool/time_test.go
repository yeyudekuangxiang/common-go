package timetool

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	DefaultTimeZone = loc
	tm, err := Parse("2006-01-02 15:04:05", "2022-01-02 15:04:05")
	assert.Equal(t, nil, err)

	assert.Equal(t, "2022-01-02 15:04:05", tm.Format("2006-01-02 15:04:05"))
}
func TestUnixMilli(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	DefaultTimeZone = loc
	msec := int64(1641107045000)
	tm := UnixMilli(msec)
	assert.Equal(t, "2022-01-02 15:04:05", tm.Format("2006-01-02 15:04:05"))
}
