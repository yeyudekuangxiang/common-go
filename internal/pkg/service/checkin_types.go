package service

type CheckinInfo struct {
	IsChecked       bool  `json:"isChecked"`
	TodayCheckIndex int   `json:"todayCheckIndex"` //今天是签到的第几天
	CheckedNumber   int   `json:"checkedNumber"`   //已签到天数
	QuickCheckin    []int `json:"quickCheckIn"`
	Rule            []int `json:"rule"`
}
