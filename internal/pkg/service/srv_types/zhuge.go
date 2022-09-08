package srv_types

import (
	"mio/internal/pkg/model/entity"
	"time"
)

type TrackPoint struct {
	OpenId       string
	PointType    entity.PointTransactionType //变动类型
	ChangeType   string                      //获取或者消耗 dec(消耗) inc(获取)
	Value        uint                        //变动数量
	IsFail       bool                        //是否失败
	FailMessage  string                      //失败信息
	AdditionInfo string                      //备注
}

type TrackOrderZhuGe struct {
	OpenId        string
	CertificateId string
	ProductItemId string
	OrderId       string
	Partnership   entity.PartnershipType
	Title         string
	CateTitle     string
	IsFail        bool   //是否失败
	FailMessage   string //失败信息
}

type TrackLoginZhuGe struct {
	OpenId      string
	IsFail      bool   //是否失败
	FailMessage string //失败信息
	Event       string
}
type TrackPoints struct {
	OpenId      string
	PointType   string //变动类型
	ChangeType  string //获取或者消耗 dec(消耗) inc(获取)
	Value       uint   //变动数量
	IsFail      bool   //是否失败
	FailMessage string //失败信息
}

type TrackBusinessPoints struct {
	Uid        string //用户编号
	Value      int
	ChangeType string    //获取或者消耗 dec(消耗) inc(获取)
	Nickname   string    //昵称
	Username   string    //姓名
	Department string    //部门
	Company    string    //公司
	ChangeTime time.Time //时间
}
type TrackBusinessCredit struct {
	Uid        string //用户编号
	Value      float64
	ChangeType string    //获取或者消耗 dec(消耗) inc(获取)
	Nickname   string    //昵称
	Username   string    //姓名
	Department string    //部门
	Company    string    //公司
	ChangeTime time.Time //时间
}
