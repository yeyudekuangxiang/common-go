package business

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/business"
	brepo "mio/internal/pkg/repository/business"
)

var DefaultPointTransactionService = PointTransactionService{repo: brepo.DefaultPointTransactionRepository}

type PointTransactionService struct {
	repo brepo.PointTransactionRepository
}

func (srv PointTransactionService) GetListBy(param GetPointTransactionListParam) []business.PointTransaction {
	return srv.repo.GetListBy(brepo.GetPointTransactionListBy{
		UserId:    param.UserId,
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		OrderBy:   param.OrderBy,
		Type:      param.Type,
	})
}

func (srv PointTransactionService) GetPointTransactionInfoList(param GetPointTransactionInfoListParam) []PointTransactionInfo {
	ptList := srv.GetListBy(GetPointTransactionListParam{
		UserId:    param.UserId,
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		OrderBy:   entity.OrderByList{business.OrderByPointTranCTDESC},
	})

	infoList := make([]PointTransactionInfo, 0)
	for _, pt := range ptList {
		infoList = append(infoList, PointTransactionInfo{
			ID:       pt.ID,
			Type:     pt.Type,
			TypeText: pt.Type.Text(),
			TimeStr:  pt.CreatedAt.Format("01.02 15:04:05"),
			Value:    pt.Value,
		})
	}
	return infoList
}
