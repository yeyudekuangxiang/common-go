package business

import (
	"github.com/pkg/errors"
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
)

var DefaultCompanyCarbonSceneService = CompanyCarbonSceneService{repo: rbusiness.DefaultCompanyCarbonSceneRepository}

type CompanyCarbonSceneService struct {
	repo rbusiness.CompanyCarbonSceneRepository
}

func (srv CompanyCarbonSceneService) FindCompanyScene(param FindCompanyCarbonSceneParam) (*ebusiness.CompanyCarbonScene, error) {
	scene := srv.repo.FindCompanyScene(rbusiness.FindCompanyCarbonSceneBy{
		CompanyId:     param.CompanyId,
		CarbonSceneId: param.CarbonSceneId,
	})
	return &scene, nil
}

//FindCompanySceneSetting 查询公司低碳场景积分和碳积分限制
func (srv CompanyCarbonSceneService) FindCompanySceneSetting(companyId int, carbonType ebusiness.CarbonType) (*CompanySceneSetting, error) {
	//需要方法-查询用户信息
	userInfo := ebusiness.User{}

	carbonScene, err := DefaultCarbonSceneService.FindScene(carbonType)
	if err != nil {
		return nil, err
	}
	if carbonScene.ID == 0 {
		return nil, errors.New("未查询到此低碳场景")
	}
	companyCarbonScene, err := DefaultCompanyCarbonSceneService.FindCompanyScene(FindCompanyCarbonSceneParam{
		CompanyId: userInfo.BCompanyId,
	})
	if err != nil {
		return nil, err
	}
	if companyCarbonScene.ID == 0 {
		return nil, errors.New("未查询到此低碳场景")
	}

	return &CompanySceneSetting{
		PointSetting:    companyCarbonScene.PointSetting,
		MaxPoint:        companyCarbonScene.MaxPoint,
		MaxCount:        companyCarbonScene.MaxCount,
		MaxCarbonCredit: carbonType.MaxDayCarbonCredit(),
	}, nil
}
