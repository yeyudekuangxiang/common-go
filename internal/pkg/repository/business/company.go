package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultCompanyRepository = CompanyRepository{DB: app.DB}

type CompanyRepository struct {
	DB *gorm.DB
}

type ICompanyRepository interface {
	Save(Company *business.Company) error
}

func (u CompanyRepository) Save(Company *business.Company) error {
	return app.DB.Save(Company).Error
}

func (u CompanyRepository) GetCompanyById(id int) business.Company {
	Company := business.Company{}
	db := app.DB.Model(Company)
	if err := db.First(&Company, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return Company
}
func (u CompanyRepository) GetCompanyPageList(by GetCompanyPageListBy) ([]business.Company, int64, error) {
	list := make([]business.Company, 0)
	var total int64
	err := u.DB.Model(business.Company{}).Offset(by.Offset).Count(&total).Limit(by.Limit).Find(&list).Error
	return list, total, err
}
