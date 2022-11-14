package common

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"time"
)

func NewCityService(ctx *context.MioContext) CityService {
	return CityService{ctx: ctx, repo: repository.NewCityRepository(ctx)}
}

type CityService struct {
	ctx  *context.MioContext
	repo repository.CityRepository
}

//  添加发放碳量记录并且更新用户剩余碳量

func (srv CityService) Create(dto CreateCityParams) (*entity.City, error) {
	//入库
	cityDo := entity.City{
		CityCode:  dto.CityCode,
		Name:      dto.Name,
		PidCode:   dto.PidCode,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now()}
	err := srv.repo.Create(&cityDo)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (srv CityService) GetByCityCode(dto GetByCityCodeParams) (entity.City, error) {
	ret, err := srv.repo.GetByCityCode(repotypes.GetCityByCode{
		CityCode: dto.CityCode,
	})
	if err != nil {
		return entity.City{}, err
	}
	return ret, nil
}

func (srv CityService) GetCity(params GetByCityCodeParams) ([]entity.City, error) {
	list, err := srv.repo.GetList(repotypes.GetCityListDO{
		PidCode: params.CityCode,
	})

	if err != nil {
		return []entity.City{}, err
	}

	return list, nil
}
