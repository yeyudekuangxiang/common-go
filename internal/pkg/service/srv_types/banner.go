package srv_types

import "mio/internal/pkg/model/entity"

type GetBannerListDTO struct {
	Scene  entity.BannerScene
	Type   entity.BannerType
	Status entity.BannerStatus
}

type CreateBannerDTO struct {
	Scene    entity.BannerScene
	Type     entity.BannerType
	Status   entity.BannerStatus
	Name     string
	ImageUrl string
	AppId    string
	Sort     int
	Redirect string
}

type UpdateBannerDTO struct {
	Id       int64
	Scene    entity.BannerScene
	Type     entity.BannerType
	Status   entity.BannerStatus
	Name     string
	ImageUrl string
	AppId    string
	Sort     int
	Redirect string
}

type GetPageBannerDTO struct {
	Name    string             `json:"name" form:"name" binding:"" alias:"banner名称"`
	Scene   string             `json:"scene" form:"scene" binding:"" alias:"轮播图场景"`
	Status  int8               `json:"status" form:"status" binding:"" alias:"上线和下线状态"`
	Offset  int                `json:"offset"`
	Limit   int                `json:"limit"` //limit为0时不限制数量
	OrderBy entity.OrderByList `json:"orderBy"`
}
