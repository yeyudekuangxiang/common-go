package repo_types

import "mio/internal/pkg/model/entity"

type GetPageFullOrderDO struct {
	Openid      string
	OrderSource entity.OrderSource
	Limit       int
	Offset      int
}
