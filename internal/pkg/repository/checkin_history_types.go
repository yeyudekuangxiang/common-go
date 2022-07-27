package repository

import "mio/internal/pkg/model/entity"

type FindCheckinHistoryBy struct {
	OpenId  string
	OrderBy entity.OrderByList
}
