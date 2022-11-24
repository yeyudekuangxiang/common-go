package timetool

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStartOfDay(t *testing.T) {
	tm1, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 00:00:00")
	tm2, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 00:00:59")
	tm3, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 00:10:59")
	tm4, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 17:53:59")
	tm5, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 23:59:59")
	assert.Equal(t, "2022-11-15 00:00:00", StartOfDay(tm1).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-11-15 00:00:00", StartOfDay(tm2).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-11-15 00:00:00", StartOfDay(tm3).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-11-15 00:00:00", StartOfDay(tm4).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-11-15 00:00:00", StartOfDay(tm5).Format("2006-01-02 15:04:05"))
}
func TestEndOfDay(t *testing.T) {
	tm1, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 00:00:00")
	tm2, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 00:00:59")
	tm3, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 00:10:59")
	tm4, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 17:53:59")
	tm5, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 23:59:59")
	assert.Equal(t, "2022-11-15 23:59:59", EndOfDay(tm1).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-11-15 23:59:59", EndOfDay(tm2).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-11-15 23:59:59", EndOfDay(tm3).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-11-15 23:59:59", EndOfDay(tm4).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-11-15 23:59:59", EndOfDay(tm5).Format("2006-01-02 15:04:05"))
}
func TestStartOfMonth(t *testing.T) {
	tm1, _ := time.Parse("2006-01-02 15:04:05", "2022-01-15 00:00:00")
	tm2, _ := time.Parse("2006-01-02 15:04:05", "2022-03-15 00:00:59")
	tm3, _ := time.Parse("2006-01-02 15:04:05", "2022-05-15 00:10:59")
	tm4, _ := time.Parse("2006-01-02 15:04:05", "2022-07-15 17:53:59")
	tm5, _ := time.Parse("2006-01-02 15:04:05", "2022-09-15 23:59:59")
	tm6, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 23:59:59")
	tm7, _ := time.Parse("2006-01-02 15:04:05", "2022-12-15 23:59:59")
	assert.Equal(t, "2022-01-15 00:00:00", StartOfDay(tm1).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-03-15 00:00:00", StartOfDay(tm2).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-05-15 00:00:00", StartOfDay(tm3).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-07-15 00:00:00", StartOfDay(tm4).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-09-15 00:00:00", StartOfDay(tm5).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-11-15 00:00:00", StartOfDay(tm6).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-12-15 00:00:00", StartOfDay(tm7).Format("2006-01-02 15:04:05"))
}
func TestEndOfMonth(t *testing.T) {
	tm1, _ := time.Parse("2006-01-02 15:04:05", "2022-01-15 00:00:00")
	tm2, _ := time.Parse("2006-01-02 15:04:05", "2022-03-15 00:00:59")
	tm3, _ := time.Parse("2006-01-02 15:04:05", "2022-05-15 00:10:59")
	tm4, _ := time.Parse("2006-01-02 15:04:05", "2022-07-15 17:53:59")
	tm5, _ := time.Parse("2006-01-02 15:04:05", "2022-09-15 23:59:59")
	tm6, _ := time.Parse("2006-01-02 15:04:05", "2022-11-15 23:59:59")
	tm7, _ := time.Parse("2006-01-02 15:04:05", "2022-12-15 23:59:59")
	assert.Equal(t, "2022-01-15 23:59:59", EndOfDay(tm1).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-03-15 23:59:59", EndOfDay(tm2).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-05-15 23:59:59", EndOfDay(tm3).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-07-15 23:59:59", EndOfDay(tm4).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-09-15 23:59:59", EndOfDay(tm5).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-11-15 23:59:59", EndOfDay(tm6).Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2022-12-15 23:59:59", EndOfDay(tm7).Format("2006-01-02 15:04:05"))
}
