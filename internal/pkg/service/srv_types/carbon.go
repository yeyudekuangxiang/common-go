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
}
type IncUserCarbonDTO struct {
	OpenId       string                       `binding:"required"`
	Type         entity.CarbonTransactionType `binding:"required"`
	BizId        string
	ChangePoint  float64
	AdminId      int
	Note         string
	AdditionInfo string
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
