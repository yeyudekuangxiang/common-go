package business

import (
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
	"mio/pkg/errno"
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
	carbonScene, err := DefaultCarbonSceneService.FindScene(carbonType)
	if err != nil {
		return nil, err
	}
	if carbonScene.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未查询到此低碳场景")
	}
	companyCarbonScene, err := DefaultCompanyCarbonSceneService.FindCompanyScene(FindCompanyCarbonSceneParam{
		CompanyId:     companyId,
		CarbonSceneId: carbonScene.ID,
	})
	if err != nil {
		return nil, err
	}
	if companyCarbonScene.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未查询到此低碳场景")
	}

	return &CompanySceneSetting{
		PointRateSetting: companyCarbonScene.PointRateSetting,
		MaxCount:         companyCarbonScene.MaxCount,
	}, nil
}

func (srv CompanyCarbonSceneService) GetBusinessCompanyCarbonSceneListBy(param rbusiness.GetCompanyCarbonSceneListBy) []ebusiness.CarbonScene {
	companyCarbonSceneList := srv.repo.GetCompanyCarbonSceneListBy(param)
	var ids []int
	for _, v := range companyCarbonSceneList {
		ids = append(ids, v.CarbonSceneId)
	}
	businessCompanyCarbonSceneList := DefaultCarbonSceneService.GetBusinessCarbonSceneListBy(rbusiness.GetCarbonSceneListBy{Ids: ids})
	return businessCompanyCarbonSceneList
}
