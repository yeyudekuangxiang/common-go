package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultUserRepository = UserRepository{DB: app.BusinessDB}

type UserRepository struct {
	DB *gorm.DB
}

type IUserRepository interface {
	Save(user *business.User) error
}

func (u UserRepository) Save(user *business.User) error {
	return u.DB.Save(user).Error
}

func (u UserRepository) GetUserBy(by GetUserBy) business.User {
	user := business.User{}
	db := u.DB.Model(user)

	if by.Uid != "" {
		db.Where("uid = ?", by.Uid)
	}
	if by.ID > 0 {
		db.Where("id = ?", by.ID)
	}
	if by.Mobile != "" {
		db.Where("mobile = ?", by.Mobile)
	}

	if err := db.First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return user
}

func (u UserRepository) GetUserListBy(by GetUserListBy) []business.User {
	list := make([]business.User, 0)
	db := u.DB.Model(business.User{})
	if len(by.Ids) > 0 {
		db.Where("id in (?)", by.Ids)
	}
	if by.CId > 0 {
		db.Where("b_company_id = ?", by.CId)
	}

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}

	return list
}
