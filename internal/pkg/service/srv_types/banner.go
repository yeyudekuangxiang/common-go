package srv_types

import "mio/internal/pkg/model/entity"

type GetBannerListDTO struct {
	Scene  entity.BannerScene
	Type   entity.BannerType
	Status entity.BannerStatus
}
