package activity

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	entity "mio/internal/pkg/model/entity/activity"
	"mio/pkg/errno"
)

func NewActivitySubjectModel(ctx *context.MioContext) SubjectModel {
	return defaultSubjectModel{ctx: ctx}
}

type SubjectModel interface {
	Save(subject *entity.Subject) error
	FindOneActive(params FindOneActiveSubjectParams) (*entity.Subject, error)
}

type defaultSubjectModel struct {
	ctx *context.MioContext
}

func (m defaultSubjectModel) Save(subject *entity.Subject) error {
	return m.ctx.DB.Model(&entity.Subject{}).Save(subject).Error
}

func (m defaultSubjectModel) FindOneActive(params FindOneActiveSubjectParams) (*entity.Subject, error) {
	var resp entity.Subject
	query := m.ctx.DB.Model(&entity.Subject{})
	if params.ActivityId != 0 {
		query = query.Where("activity_id = ?", params.ActivityId)
	}
	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	}
	if !params.StartTime.IsZero() {
		query = query.Where("start_time > ?", params.StartTime)
	}
	if !params.EndTime.IsZero() {
		query = query.Where("end_time < ?", params.EndTime)
	}
	if params.Title != "" {
		query = query.Where("title like ?", params.Title+"%")
	}
	if params.Updater != "" {
		query = query.Where("updater = ?", params.Updater)
	}
	if params.Creator != "" {
		query = query.Where("creator = ?", params.Creator)
	}
	if params.UperRiskLimit != 0 {
		query = query.Where("uper_risk_limit = ?", params.UperRiskLimit)
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
