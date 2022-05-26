package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultUserRepository IUserRepository = NewUserRepository()

type IUserRepository interface {
	Save(user *entity.User) error
	// GetUserById 根据用id获取用户信息
	GetUserById(int64) entity.User
	GetUserBy(by GetUserBy) entity.User
	GetUserListBy(by GetUserListBy) []entity.User
	GetGuid(unionId string) string
	GetUserPageListBy(by GetUserPageListBy) ([]entity.User, int64)
}

func NewUserRepository() UserRepository {
	return UserRepository{}
}

type UserRepository struct {
}

func (u UserRepository) Save(user *entity.User) error {
	return app.DB.Save(user).Error
}
func (u UserRepository) GetUserById(id int64) entity.User {
	var user entity.User
	if err := app.DB.First(&user, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return user
}
func (u UserRepository) GetUserBy(by GetUserBy) entity.User {
	user := entity.User{}
	db := app.DB.Model(user)

	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if by.Source != "" {
		db.Where("source = ?", by.Source)
	}
	if by.Mobile != "" {
		db.Where("phone_number = ?", by.Mobile)
	}
	if by.LikeMobile != "" {
		db.Where("phone_number like ?", "%"+by.LikeMobile+"%")
	}
	if by.UnionId != "" {
		db.Where("unionid = ?", by.UnionId)
	}

	if err := db.First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return user
}

func (u UserRepository) GetUserListBy(by GetUserListBy) []entity.User {
	list := make([]entity.User, 0)
	db := app.DB.Model(entity.User{})

	if by.Mobile != "" {
		db.Where("phone_number = ?", by.Mobile)
	}
	if len(by.UserIds) > 0 {
		db.Where("id in (?)", by.UserIds)
	}
	if len(by.Mobiles) > 0 {
		db.Where("phone_number in (?)", by.Mobiles)
	}
	if by.Source != "" {
		db.Where("source = ?", by.Source)
	}
	if by.Nickname != "" {
		db.Where("nick_name like ?", "%"+by.Nickname+"%")
	}
	if by.LikeMobile != "" {
		db.Where("phone_number like ?", "%"+by.LikeMobile+"%")
	}
	if by.UserId != 0 {
		db.Where("id = ?", by.UserId)
	}
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if !by.StartTime.IsZero() {
		db.Where("time >= ?", by.StartTime)
	}
	if !by.EndTime.IsZero() {
		db.Where("time <= ?", by.EndTime)
	}

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
func (u UserRepository) GetGuid(unionId string) string {
	if unionId == "" {
		return ""
	}
	user := entity.User{}
	err := app.DB.Where("unionid = ? and guid <> ''", unionId).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return user.GUID
}

func (u UserRepository) GetUserPageListBy(bp GetUserPageListBy) ([]entity.User, int64) {
	list := make([]entity.User, 0)
	var count int64
	db := app.DB.Model(entity.User{})
	by := bp.User
	if by.Mobile != "" {
		db.Where("phone_number = ?", by.Mobile)
	}
	if len(by.UserIds) > 0 {
		db.Where("id in (?)", by.UserIds)
	}
	if len(by.Mobiles) > 0 {
		db.Where("phone_number in (?)", by.Mobiles)
	}
	if by.Source != "" {
		db.Where("source = ?", by.Source)
	}
	if by.Nickname != "" {
		db.Where("nick_name like ?", "%"+by.Nickname+"%")
	}
	if by.LikeMobile != "" {
		db.Where("phone_number like ?", "%"+by.LikeMobile+"%")
	}
	if by.UserId != 0 {
		db.Where("id = ?", by.UserId)
	}
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if !by.StartTime.IsZero() {
		db.Where("time >= ?", by.StartTime)
	}
	if !by.EndTime.IsZero() {
		db.Where("time <= ?", by.EndTime)
	}

	if err := db.Find(&list).Limit(bp.Limit).Offset(bp.Offset).Order(bp.OrderBy).Count(&count).Error; err != nil {
		panic(err)
	}
	return list, count
}
