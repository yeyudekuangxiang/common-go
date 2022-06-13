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

func (u CompanyService) GetCompanyById(id int) (*business.Company, error) {
	if id == 0 {
		return &business.Company{}, nil
	}
	company := u.repo.GetCompanyById(id)
	return &company, nil
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
