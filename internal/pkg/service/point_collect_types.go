package service

import "mio/internal/pkg/model/entity"

type CreateHistoryParam struct {
	OpenId          string
	TransactionType entity.PointTransactionType
	Info            string
}
