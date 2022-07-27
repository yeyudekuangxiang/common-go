package service

import "mio/internal/pkg/model/entity"

type FindCheckinHistoryParam struct {
	OpenId  string
	OrderBy entity.OrderByList
}
