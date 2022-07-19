package repotypes

import "mio/internal/pkg/model/entity"

type GetBannerOneDO struct {
	ID int64
}

type GetBannerListDO struct {
	Scene   entity.BannerScene
	Type    entity.BannerType
	Status  entity.BannerStatus
	OrderBy entity.OrderByList
}

type GetBannerPageDO struct {
	Scene   entity.BannerScene
	Type    entity.BannerType
	OrderBy entity.OrderByList
	Offset  int
	Limit   int
}
