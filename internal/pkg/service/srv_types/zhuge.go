package srv_types

import "mio/internal/pkg/model/entity"

type TrackPoint struct {
	OpenId      string
	PointType   entity.PointTransactionType //变动类型
	ChangeType  string                      //获取或者消耗 dec(消耗) inc(获取)
	Value       uint                        //变动数量
	IsFail      bool                        //是否失败
	FailMessage string                      //失败信息
}

type TrackPoints struct {
	OpenId      string
	PointType   string //变动类型
	ChangeType  string //获取或者消耗 dec(消耗) inc(获取)
	Value       uint   //变动数量
	IsFail      bool   //是否失败
	FailMessage string //失败信息
}
