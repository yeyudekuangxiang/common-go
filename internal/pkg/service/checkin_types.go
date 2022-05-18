package service

type CheckinInfo struct {
	IsChecked     bool  `json:"isChecked"`
	NthDayToCheck int   `json:"nthDayToCheck"`
	QuickCheckin  []int `json:"quickCheckIn"`
	Rule          []int `json:"rule"`
}
