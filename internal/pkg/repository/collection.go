package repository

import (
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"time"
)

type (
	CollectionModel interface {
		FindOne(id int64) (*entity.Collection, error)
		Insert(data *entity.Collection) (*entity.Collection, error)
		Delete(id int64) error
		Update(data *entity.Collection) error
		FindAllByOpenId(openId string) ([]*entity.Collection, error)
		FindAllByTime(startTime, endTime time.Time) ([]*entity.Collection, error)
	}

	defaultCollectionModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultCollectionModel) FindOne(id int64) (*entity.Collection, error) {
	//TODO implement me
	panic("implement me")
}

func (d defaultCollectionModel) Insert(data *entity.Collection) (*entity.Collection, error) {
	//TODO implement me
	panic("implement me")
}

func (d defaultCollectionModel) Delete(id int64) error {
	//TODO implement me
	panic("implement me")
}

func (d defaultCollectionModel) Update(data *entity.Collection) error {
	//TODO implement me
	panic("implement me")
}

func (d defaultCollectionModel) FindAllByOpenId(openId string) ([]*entity.Collection, error) {
	//TODO implement me
	panic("implement me")
}

func (d defaultCollectionModel) FindAllByTime(startTime, endTime time.Time) ([]*entity.Collection, error) {
	//TODO implement me
	panic("implement me")
}

func NewCollectionRepository() CollectionModel {
	return &defaultCollectionModel{
		ctx: mioContext.NewMioContext(),
	}
}
