package entity

import "mio/internal/pkg/model"

type Step struct {
	ID                 int64
	OpenId             string     `gorm:"column:openid"`
	Total              int64      //步行总数量
	LastCheckTime      model.Time //领取积分时stepHistory最后一条记录的recorded_time(记录最后一次领积分的日期 如果是今天计算积分时需要减去已领取的步数 如果不是今天就可以直接计算积分)
	LastCheckCount     int        //领取积分时stepHistory最后一条数量的count
	LastSumHistoryTime model.Time //最有一次被计算进总步数的history的recorded_time(用户可能之前没更新步数 然后一次会更新多天的步数 所以需要记录上次计算总数的开始日期)
	LastSumHistoryNum  int        //最有一次被计算进总步数的history的count
}

func (Step) TableName() string {
	return "step"
}
