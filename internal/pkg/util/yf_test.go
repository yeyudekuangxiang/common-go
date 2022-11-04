package util

import (
	"fmt"
	"math"
	"testing"
)

const (
	AVG_METRIC_AGE float64 = 30.0
	DECAY          float64 = 2 / (float64(AVG_METRIC_AGE) + 1)
)

type SimpleEWMA struct {
	// 当前平均值。在用Add()添加后，这个值会更新所有数值的平均值。
	value float64
}

// 添加并更新滑动平均值
func (e *SimpleEWMA) Add(value float64) {
	if e.value == 0 { // this is a proxy for "uninitialized"
		e.value = value
	} else {
		e.value = (value * DECAY) + (e.value * (1 - DECAY))
	}
}

// 获取当前滑动平均值
func (e *SimpleEWMA) Value() float64 {
	return e.value
}

// 设置 ewma 值
func (e *SimpleEWMA) Set(value float64) {
	e.value = value
}

const testMargin = 0.00000001

var samples = [100]float64{
	4599, 5711, 4746, 4621, 5037, 4218, 4925, 4281, 5207, 5203, 5594, 5149,
	4948, 4994, 6056, 4417, 4973, 4714, 4964, 5280, 5074, 4913, 4119, 4522,
	4631, 4341, 4909, 4750, 4663, 5167, 3683, 4964, 5151, 4892, 4171, 5097,
	3546, 4144, 4551, 6557, 4234, 5026, 5220, 4144, 5547, 4747, 4732, 5327,
	5442, 4176, 4907, 3570, 4684, 4161, 5206, 4952, 4317, 4819, 4668, 4603,
	4885, 4645, 4401, 4362, 5035, 3954, 4738, 4545, 5433, 6326, 5927, 4983,
	5364, 4598, 5071, 5231, 5250, 4621, 4269, 3953, 3308, 3623, 5264, 5322,
	5395, 4753, 4936, 5315, 5243, 5060, 4989, 4921, 4480, 3426, 3687, 4220,
	3197, 5139, 6101, 5279,
}

func withinMargin(a, b float64) bool {
	return math.Abs(a-b) <= testMargin
}

func TestSimpleEWMA(t *testing.T) {
	var e SimpleEWMA
	for _, f := range samples {
		e.Add(f)
	}
	fmt.Println(e.Value())
	if !withinMargin(e.Value(), 4734.500946466118) {
		t.Errorf("e.Value() is %v, wanted %v", e.Value(), 4734.500946466118)
	}
	e.Set(1.0)
	if e.Value() != 1.0 {
		t.Errorf("e.Value() is %v", e.Value())
	}
}
