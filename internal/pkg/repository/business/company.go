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
