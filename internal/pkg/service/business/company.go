package business

import (
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model/entity/business"
	brepo "mio/internal/pkg/repository/business"
)

var DefaultCompanyService = CompanyService{repo: brepo.DefaultCompanyRepository}

type CompanyService struct {
	repo brepo.CompanyRepository
}

func (u CompanyService) GetCompanyById(id int) *business.Company {
	if id == 0 {
		return &business.Company{}
	}
	company := u.repo.GetCompanyById(id)
	return &company
}

// GetCompanyCarbon 当月企业碳减排综合
func (u CompanyService) GetCompanyCarbon(cid int) decimal.Decimal {
	a := decimal.Decimal{}
	list, err := DefaultCarbonRankService.DepartmentRankListByCid(cid)
	if err != nil || len(list) == 0 || list[0].ID == 0 {
		return a
	}

	for _, v := range list {
		a = a.Add(v.Value)
	}
	return a
}
func (u CompanyService) GetCompanyPageList(param GetCompanyPageListParam) ([]business.Company, int64, error) {
	return u.repo.GetCompanyPageList(brepo.GetCompanyPageListBy{
		Limit:  param.Limit,
		Offset: param.Offset,
	})
}
