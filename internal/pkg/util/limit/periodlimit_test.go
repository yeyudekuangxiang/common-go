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

type MovingAverage struct {
	index   int   // 当前环形数组的位置
	count   int   // 数组大小
	sum     int   // 数据总量
	buckets []int // 环形数组
}

/** Initialize your data structure here. */
func Constructor(size int) MovingAverage {
	return MovingAverage{index: size - 1, buckets: make([]int, size)}
}

func (ma *MovingAverage) Next(val int) float64 {
	ma.sum += val
	ma.index = (ma.index + 1) % len(ma.buckets) // 循环数组索引
	if ma.count < len(ma.buckets) {
		ma.count++
		ma.buckets[ma.index] = val
	} else {
		ma.sum -= ma.buckets[ma.index] // 减去旧数据
		ma.buckets[ma.index] = val     // 替换旧数据
	}
	return float64(ma.sum) / float64(ma.count)
}

func Test_Demo(t *testing.T) {
	//ma := Constructor(3)
	//fmt.Println(ma.Next(1))  // 返回 1.0 = 1 / 1
	//fmt.Println(ma.Next(10)) // 返回 5.5 = (1 + 10) / 2
	//fmt.Println(ma.Next(3))  // 返回 4.66667 = (1 + 10 + 3) / 3
	//fmt.Println(ma.Next(5))  // 返回 6.0 = (10 + 3 + 5) / 3
	d := time.Now().AddDate(0, 1, -time.Now().Day()).Day()
	fmt.Println(d) // 返回 6.0 = (10 + 3 + 5) / 3
	d2 := time.Now().AddDate(0, 0, 0)

	fmt.Println(d2)
	fmt.Println(time.Now().Day())
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	fmt.Println(firstDay)
	fmt.Println(now)
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Nanosecond)
	fmt.Println(lastDay.Unix() - time.Now().Unix())
	fmt.Println(lastDay.Day())
}
