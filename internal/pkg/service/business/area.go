package business

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
	"mio/internal/pkg/util"
	"sort"
	"strings"
)

type AreaService struct {
	ctx  *context.MioContext
	repo *rbusiness.AreaRepository
}

func NewAreaService(ctx *context.MioContext) *AreaService {
	return &AreaService{ctx: ctx, repo: rbusiness.NewAreaRepository(ctx)}
}

func (srv AreaService) List(dto AreaListDTO) ([]ebusiness.Area, error) {
	po := rbusiness.AreaListPO{}
	if err := util.MapTo(dto, &po); err != nil {
		return nil, err
	}

	return srv.repo.List(po)
}

func (srv AreaService) GroupCityProvinceList(dto CityProvinceListDTO) ([]GroupCityProvince, error) {
	alDto := AreaListDTO{}
	if err := util.MapTo(dto, &alDto); err != nil {
		return nil, err
	}
	alDto.Level = ebusiness.AreaCity

	cityList, err := srv.List(alDto)
	if err != nil {
		return nil, err
	}
	parentCityIds := make([]int64, 0)

	for _, city := range cityList {
		parentCityIds = append(parentCityIds, int64(city.ParentCityID))
	}
	provinceList, err := srv.repo.List(rbusiness.AreaListPO{
		CityIds: parentCityIds,
		Level:   ebusiness.AreaProvince,
	})
	if err != nil {
		return nil, err
	}
	provinceMap := make(map[model.LongID]ebusiness.Area)
	for _, province := range provinceList {
		provinceMap[province.CityID] = province
	}

	ctMap := make(map[string][]CityProvince)

	for _, city := range cityList {
		c := ShortArea{}
		if err := util.MapTo(city, &c); err != nil {
			return nil, err
		}
		p := ShortArea{}
		if err := util.MapTo(provinceMap[city.ParentCityID], &p); err != nil {
			return nil, err
		}

		letter := strings.ToUpper(city.Py[:1])
		ctMap[letter] = append(ctMap[letter], CityProvince{
			Province: p,
			City:     c,
		})
	}

	gcpList := make([]GroupCityProvince, 0)
	for letter, g := range ctMap {
		gcpList = append(gcpList, GroupCityProvince{
			Letter: letter,
			Items:  g,
		})
	}

	sort.Slice(gcpList, func(i, j int) bool {
		return gcpList[i].Letter < gcpList[j].Letter
	})

	return gcpList, nil
}
