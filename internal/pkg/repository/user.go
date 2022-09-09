package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultUserRepository = NewUserRepository()

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

func (u UserRepository) GetUserIdentifyInfo(openid string) (info GetUserIdentifyInfoBy, exist bool, err error) {
	var ret GetUserIdentifyInfoBy
	var user entity.User
	db := app.DB.Model(user).Joins("left join user_channel on user_channel.cid  = \"user\".channel_id")
	db.Joins("left join city on  city.city_code = \"user\".city_code")
	db.Joins("left join point on point.openid = \"user\".openid")
	db.Joins("left join invite on invite.new_user_openid = \"user\".openid")
	db.Joins("left join user_channel_type on user_channel_type.id = user_channel.pid")
	//db.Select("user.openid,user.nick_name,user.time")
	db.Select("\"user\".openid,\"user\".nick_name,\"user\".time,\"user\".source,user_channel_type.name as channel_type_name,user_channel.name as channel_name,city.name as city_name,point.balance,invite.invited_by_openid ")
	db.Where("user.openid", openid)
	err = db.First(&ret).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ret, false, nil
		}
		return ret, false, err
	}
	return ret, true, nil
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
	if by.Risk != 0 {
		db.Where("risk = ?", by.Risk)
	}
	if by.Status != 0 {
		db.Where("status = ?", by.Status)
	}
	if by.Partners != 0 {
		db.Where("partners = ?", by.Partners)
	}
	if by.Position != "" {
		db.Where("position = ?", by.Position)
	}
	if by.Auth != 0 {
		db.Where("auth = ?", by.Auth)
	}

	if err := db.Count(&count).Limit(bp.Limit).Offset(bp.Offset).Order(bp.OrderBy).Find(&list).Error; err != nil {
		panic(err)
	}
	return list, count
}

func (u UserRepository) GetUserByID(id int64) (*entity.User, bool, error) {
	var user entity.User
	if err := app.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &user, true, nil
}
func (u UserRepository) GetUser(by GetUserBy) (*entity.User, bool, error) {
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
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &user, true, nil
}
