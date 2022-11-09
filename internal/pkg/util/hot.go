package util

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"time"
)

type Hot struct{}

func NewHot() *Hot {
	return &Hot{}
}

const (
	col = 0.223
)

/*
牛顿冷却 : 本期得分 = 上一期得分 x exp(-(冷却系数) x 间隔的小时数)
seeCount : 累计查看数 每个1分
likeCount : 累计点赞数 每个2分
commentCount: 累计评论数 每个3分
exp : 欧拉数
a ： 冷却系数 如第一天100分 第二天自然冷却到80分； 反推得到系数为 0.223
t : 出始时间 - 当前时间
*/
func (n *Hot) Hot(views, likes, comments int64, createdTime time.Time) float64 {
	//本期热度 = (seeCount * 1 + likeCount * 2 + commentCount * 3) * exp^(-a*t)
	//t := createdTime.Sub(time.Now()).Hours()
	t := time.Now().Sub(createdTime).Hours() / 24
	exp := math.Exp(-col * t)
	fmt.Printf("exp:%f; views: %d, likes:%d, comments:%d; t:%f \n", exp, views, likes, comments, t)
	high, _ := decimal.NewFromInt(views + likes*2 + comments*3).Mul(decimal.NewFromFloat(exp)).Float64()
	return high
}

/*

log10(Qviews)*4	: 浏览次数
Qscore*Qanswers/5 : Qscore 得分 Qanswers 评论
Ascores: 回答得分
Qage: 距离问题发表的时间
Qupdated: 距离最后一个回答的时间

*/
func (n *Hot) HotV2(Qviews, Qanswers, Qscore, Ascores float64, dataAsk, dateActive time.Time) {
	Qage := time.Now().Sub(dataAsk).Seconds()
	Qage, _ = decimal.NewFromFloat(Qage / 3600).Round(1).Float64()

	Qupdate := time.Now().Sub(dateActive).Seconds()
	Qupdate, _ = decimal.NewFromFloat(Qupdate / 3600).Round(1).Float64()

	dividend := (math.Log10(Qviews) * 4) + ((Qanswers * Qscore) / 5) + Ascores
	divisor := math.Pow((Qage+1)-(Qage-Qupdate)/2, 1.5)
	res := dividend / divisor
	fmt.Println(res)
}
