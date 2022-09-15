package question

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	questionEntity "mio/internal/pkg/model/entity/question"
	"mio/internal/pkg/repository/repotypes"
)

func NewUserRepository(ctx *context.MioContext) *UserRepository {
	return &UserRepository{
		ctx: ctx,
	}
}

type UserRepository struct {
	ctx *context.MioContext
}

func (repo UserRepository) Save(transaction *questionEntity.User) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo UserRepository) Create(transaction *questionEntity.User) error {
	return repo.ctx.DB.Create(transaction).Error
}

func (repo UserRepository) FindBy(by repotypes.GetQuestUserGetById) questionEntity.User {
	record := questionEntity.User{}
	db := app.DB.Model(questionEntity.User{})
	if by.UserId > 0 {
		db.Where("user_id = ?", by.UserId)
	}
	if by.OpenId != "" {
		db.Where("third_id = ?", by.OpenId)
	}
	if err := db.First(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}
