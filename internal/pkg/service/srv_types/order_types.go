package srv_types

import "mio/internal/pkg/model/entity"

type GetPageFullOrderDTO struct {
	Openid      string
	OrderSource entity.OrderSource
	Offset      int
	Limit       int
}
