package activity

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	entity "mio/internal/pkg/model/entity/activity"
	"mio/pkg/errno"
)

func NewActivityModel(ctx *context.MioContext) ActivitiesModel {
	return defaultActivityModel{ctx: ctx}
}

type ActivitiesModel interface {
	Save(activity *entity.Activity) error
}

type defaultActivityModel struct {
	ctx *context.MioContext
}

func (m defaultActivityModel) Save(activity *entity.Activity) error {
	return m.ctx.DB.Model(&entity.Activity{}).Save(activity).Error
}

func (m defaultActivityModel) FindOneActive(params FindOneActiveActivityParams) (*entity.Activity, error) {
	var resp entity.Activity
	query := m.ctx.DB.Model(&entity.Activity{}).Preload("subject", func(db *gorm.DB) *gorm.DB {
		if params.SubjectStatus != 0 {
			db = db.Where("activity_subject.status = ?", params.SubjectStatus)
		}
		if !params.SubjectStartTime.IsZero() {
			db = db.Where("start_time > ?", params.SubjectStartTime)
		}
		if !params.SubjectEndTime.IsZero() {
			db = db.Where("end_time < ?", params.SubjectEndTime)
		}
		if params.SubjectCreator != "" {
			db = db.Where("creator = ?", params.SubjectCreator)
		}
		if params.SubjectUpdater != "" {
			db = db.Where("updater = ?", params.SubjectUpdater)
		}
		return db
	})

	if params.Id != 0 {
		query = query.Where("activity.id = ?", params.Id)
	}

	if params.Title != "" {
		query = query.Where("activity.title = ?", params.Title)
	}

	if params.Code != "" {
		query = query.Where("activity.code = ?", params.Code)
	}

	err := query.Take(&resp).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, errno.ErrRecordNotFound
	default:
		return nil, err
	}
}
