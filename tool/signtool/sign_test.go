package signtool

import "testing"

func Test_sign(t *testing.T) {
	sign := GetSign(map[string]interface{}{
		"city":      "100100",
		"endTime":   1684986902627,
		"mileage":   100,
		"orderNo":   "123456abc",
		"phone":     "18840851111",
		"startTime": 1684986902627,
		"source":    "mio",
	}, "0tlrEVZtRE", "&")
	println(sign)
}
