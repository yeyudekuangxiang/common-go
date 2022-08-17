package srv_types

import "mio/internal/pkg/model/entity"

type ChangeUserCarbonDTO struct {
	OpenId       string                       `binding:"required"`
	Type         entity.CarbonTransactionType `binding:"required"`
	ChangePoint  float64
	AdminId      int
	Note         string
	AdditionInfo string
	BizId        string
	CityCode     string
	Uid          int64
}
type IncUserCarbonDTO struct {
	OpenId       string                       `binding:"required"`
	Type         entity.CarbonTransactionType `binding:"required"`
	BizId        string
	ChangePoint  float64
	AdminId      int
	Note         string
	AdditionInfo string
	CityCode     string
	Uid          int64
}
type DecUserCarbonDTO struct {
	OpenId       string                      `binding:"required"`
	Type         entity.PointTransactionType `binding:"required"`
	ChangePoint  float64
	BizId        string
	AdminId      int
	Note         string
	AdditionInfo string
}
