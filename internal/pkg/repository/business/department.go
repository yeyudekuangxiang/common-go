package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultDepartmentRepository = DepartmentRepository{DB: app.BusinessDB}

type DepartmentRepository struct {
	DB *gorm.DB
}

type IDepartmentRepository interface {
	Save(Department *business.Department) error
}

func (u DepartmentRepository) Save(Department *business.Department) error {
	return u.DB.Save(Department).Error
}

func (u DepartmentRepository) GetDepartmentBy(by GetDepartmentBy) business.Department {
	Department := business.Department{}
	db := u.DB.Model(Department)

	if by.ID > 0 {
		db.Where("id = ?", by.ID)
	}

	if err := db.First(&Department).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return Department
}

func (u DepartmentRepository) GetDepartmentListBy(by GetDepartmentListBy) []business.Department {
	list := make([]business.Department, 0)
	db := u.DB.Model(business.Department{})
	if len(by.Ids) > 0 {
		db.Where("id in (?)", by.Ids)
	}

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}

	return list
}
