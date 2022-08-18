package srv_types

import "mio/internal/pkg/model/entity"

type ChangeUserPointDTO struct {
	OpenId       string                      `binding:"required"`
	Type         entity.PointTransactionType `binding:"required"`
	ChangePoint  int64
	AdminId      int
	Note         string
	AdditionInfo string
	BizId        string
	InviteId     int64
}
type IncUserPointDTO struct {
	OpenId       string                      `binding:"required"`
	Type         entity.PointTransactionType `binding:"required"`
	BizId        string
	ChangePoint  int64
	AdminId      int
	Note         string
	AdditionInfo string
	InviteId     int64
}
type DecUserPointDTO struct {
	OpenId       string                      `binding:"required"`
	Type         entity.PointTransactionType `binding:"required"`
	ChangePoint  int64
	BizId        string
	AdminId      int
	Note         string
	AdditionInfo string
}
