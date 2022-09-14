package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultUserSpecialRepository = NewUserSpecialRepository()

func NewUserSpecialRepository() UserSpecialRepository {
	return UserSpecialRepository{}
}

type UserSpecialRepository struct {
}

func (u UserSpecialRepository) Save(UserSpecial *entity.UserSpecial) error {
	return app.DB.Save(UserSpecial).Error
}

func (u UserSpecialRepository) GetUserSpecialByPhone(phone string) entity.UserSpecial {
	var UserSpecial entity.UserSpecial
	if err := app.DB.Where("phone = ?", phone).First(&UserSpecial).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return UserSpecial
}
