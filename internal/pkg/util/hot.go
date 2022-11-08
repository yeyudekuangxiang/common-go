package util

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"time"
)

/*

log10(Qviews)*4	: 浏览次数
Qscore*Qanswers/5 : Qscore 得分 Qanswers 评论
Ascores: 回答得分
Qage: 距离问题发表的时间
Qupdated: 距离最后一个回答的时间

*/
func hot(Qviews, Qanswers, Qscore, Ascores float64, dataAsk, dateActive time.Time) {
	Qage := time.Now().Sub(dataAsk).Seconds()
	Qage, _ = decimal.NewFromFloat(Qage / 3600).Round(1).Float64()

	Qupdate := time.Now().Sub(dateActive).Seconds()
	Qupdate, _ = decimal.NewFromFloat(Qupdate / 3600).Round(1).Float64()

	dividend := (math.Log10(Qviews) * 4) + ((Qanswers * Qscore) / 5) + Ascores
	divisor := math.Pow((Qage+1)-(Qage-Qupdate)/2, 1.5)
	res := dividend / divisor
	fmt.Println(res)
}
