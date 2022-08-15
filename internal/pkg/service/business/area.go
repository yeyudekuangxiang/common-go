package business

import (
	"mio/internal/pkg/core/context"
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
	"mio/internal/pkg/util"
	"sort"
)

type AreaService struct {
	ctx  *context.MioContext
	repo *rbusiness.AreaRepository
}

func NewAreaService(ctx *context.MioContext) *AreaService {
	return &AreaService{ctx: ctx, repo: rbusiness.NewAreaRepository(ctx)}
}

func (srv AreaService) CityList(dto CityLisDTO) ([]ebusiness.Area, error) {
	po := rbusiness.CityLisPO{}
	if err := util.MapTo(dto, &po); err != nil {
		return nil, err
	}

	return srv.repo.CityList(po)
}

func (srv AreaService) List(dto AreaListDTO) ([]ebusiness.Area, error) {
	po := rbusiness.AreaListPO{}
	if err := util.MapTo(dto, &po); err != nil {
		return nil, err
	}

	return srv.repo.List(po)
}

func (srv AreaService) CityProvinceList(dto CityLisDTO) ([]CityProvince, error) {
	cityList, err := srv.CityList(dto)
	if err != nil {
		return nil, err
	}
	cityCodes := make([]string, 0)
	for _, city := range cityList {
		cityCodes = append(cityCodes, city.PidCode)
	}

	provinceList, err := srv.repo.List(rbusiness.AreaListPO{
		CityCodes: cityCodes,
	})
	if err != nil {
		return nil, err
	}
	provinceMap := make(map[string]ebusiness.Area)
	for _, province := range provinceList {
		provinceMap[province.CityCode] = province
	}

	sort.Slice(cityList, func(i, j int) bool {
		return cityList[i].Name < cityList[j].Name
	})

	cityProvinceList := make([]CityProvince, 0)
	for _, city := range cityList {
		c := ShortArea{}
		if err := util.MapTo(city, &c); err != nil {
			return nil, err
		}
		p := ShortArea{}
		if err := util.MapTo(provinceMap[city.PidCode], &p); err != nil {
			return nil, err
		}
		cityProvinceList = append(cityProvinceList, CityProvince{
			Province: p,
			City:     c,
		})
	}
	return cityProvinceList, nil
}
