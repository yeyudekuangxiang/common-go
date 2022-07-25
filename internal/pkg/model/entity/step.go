package entity

import "mio/internal/pkg/model"

type Step struct {
	ID                 int64
	OpenId             string     `gorm:"column:openid"`
	Total              int64      //步行总数量
	LastCheckTime      model.Time //领取积分时stepHistory最后一条记录的recorded_time
	LastCheckCount     int        //领取积分时stepHistory最后一条数量的count
	LastSumHistoryTime model.Time //最有一次被计算进总步数的history的recorded_time
	LastSumHistoryNum  int        //最有一次被计算进总步数的history的count
}

func (Step) TableName() string {
	return "step"
}
