package repository

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repotypes"
)

func NewUserFriendRepository(ctx *context.MioContext) UserFriendRepository {
	return UserFriendRepository{ctx: ctx}
}

type UserFriendRepository struct {
	ctx *context.MioContext
}

func (repo UserFriendRepository) Save(transaction *entity.UserFriend) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo UserFriendRepository) Create(transaction *entity.UserFriend) error {
	return repo.ctx.DB.Create(transaction).Error
}

func (repo UserFriendRepository) GetListBy(by repotypes.GetUserFriendListBy) []entity.UserFriend {
	list := make([]entity.UserFriend, 0)
	db := repo.ctx.DB.Model(entity.UserFriend{})
	if len(by.UserIds) > 0 {
		db.Where("uid in (?)", by.UserIds)
	}
	if len(by.FUserIds) > 0 {
		db.Where("f_uid in (?)", by.FUserIds)
	}
	if by.FUid != 0 {
		db.Where("f_uid = ?", by.FUid)
	}
	if by.Uid != 0 {
		db.Where("uid = ?", by.Uid)
	}
	if by.Type != 0 {
		db.Where("type = ?", by.Type)
	}
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
