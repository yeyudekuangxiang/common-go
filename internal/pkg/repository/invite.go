package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

func NewInviteRepository(ctx *context.MioContext) *InviteRepository {
	return &InviteRepository{
		ctx: ctx,
	}
}

type InviteRepository struct {
	ctx *context.MioContext
}

func (repo InviteRepository) Save(transaction *entity.Invite) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo InviteRepository) Create(transaction *entity.Invite) error {
	return repo.ctx.DB.Create(transaction).Error
}

func (repo InviteRepository) GetInvite(openid string) *entity.Invite {
	invite := entity.Invite{}
	err := app.DB.Where("new_user_openid = ? and invited_by_openid <> '' and is_reward = 1", openid).First(&invite).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	if invite.ID != 0 {
		return &invite
	}
	return &entity.Invite{}
}

func (repo InviteRepository) UpdateIsReward(id int64) error {
	var result entity.Invite
	err := repo.ctx.DB.Where("id = ?", id).First(&result).Error
	if err != nil {
		return err
	}
	result.IsReward = 0
	return repo.ctx.DB.Save(&result).Error
}
