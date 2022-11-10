package repotypes

import "mio/internal/pkg/model/entity"

type GetIndexIconOneDO struct {
	ID int64
}

type GetIndexIconExistDO struct {
	Name     string
	ImageUrl string
	NotId    int64
}
type GetIndexIconDO struct {
	Scene   entity.BannerScene
	Type    entity.BannerType
	Status  entity.BannerStatus
	OrderBy entity.OrderByList
}

type DeleteIndexIconDO struct {
	Id     int64
	Status int8
}
