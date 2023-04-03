package community

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"time"
)

type (
	ActivitiesSignupModel interface {
		FindAllAPISignup(params FindAllActivitiesSignupParams) ([]*entity.APIActivitiesSignup, int64, error)
		FindOneAPISignup(params FindOneActivitiesSignupParams) (*entity.APIActivitiesSignup, error)
		FindSignupList(params FindAllActivitiesSignupParams) ([]*entity.APISignupList, int64, error)
		FindOne(params FindOneActivitiesSignupParams) (*entity.CommunityActivitiesSignup, error)
		FindAll(params FindAllActivitiesSignupParams) ([]*entity.CommunityActivitiesSignup, int64, error)
		CancelSignup(signup *entity.CommunityActivitiesSignup) error
		Delete(id int64) error
		Update(signup *entity.CommunityActivitiesSignup) error
		Create(signup *entity.CommunityActivitiesSignup) error
		FindListCount(params FindListCountParams) ([]*entity.APIListCount, error)
	}

	defaultCommunityActivitiesSignupModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultCommunityActivitiesSignupModel) FindListCount(params FindListCountParams) ([]*entity.APIListCount, error) {
	var resp []*entity.APIListCount

	err := d.ctx.DB.Model(&entity.CommunityActivitiesSignup{}).
		Select([]string{"topic_id", "COUNT(id) as num_of_signup"}).
		Where("topic_id in ?", params.TopicIds).
		Where("signup_status = 1").
		Group("topic_id").
		Find(&resp).Error

	if err != nil {
		return nil, err
	}
	return resp, err
}

func (d defaultCommunityActivitiesSignupModel) FindSignupList(params FindAllActivitiesSignupParams) ([]*entity.APISignupList, int64, error) {
	list := make([]*entity.APISignupList, 0)
	var total int64
	db := d.ctx.DB.Model(&entity.CommunityActivitiesSignup{}).
		Preload("User")

	if params.TopicId != 0 {
		db.Where("topic_id = ?", params.TopicId)
	}
	if params.UserId != 0 {
		db.Where("user_id = ?", params.UserId)
	}
	if params.City != "" {
		db.Where("city = ?", params.City)
	}
	if params.Age != 0 {
		db.Where("age = ?", params.Age)
	}
	if params.Gender != 0 {
		db.Where("gender = ?", params.Gender)
	}
	if params.Phone != "" {
		db.Where("phone = ?", params.Phone)
	}
	if params.RealName != "" {
		db.Where("real_name = ?", params.RealName)
	}

	if params.Wechat != "" {
		db.Where("wechat = ?", params.Wechat)
	}

	err := db.Where("signup_status = 1").Count(&total).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (d defaultCommunityActivitiesSignupModel) FindAll(params FindAllActivitiesSignupParams) ([]*entity.CommunityActivitiesSignup, int64, error) {
	list := make([]*entity.CommunityActivitiesSignup, 0)
	var total int64
	db := d.ctx.DB.Model(&entity.CommunityActivitiesSignup{})

	if params.TopicId != 0 {
		db.Where("topic_id = ?", params.TopicId)
	}
	if params.UserId != 0 {
		db.Where("user_id = ?", params.UserId)
	}
	if params.City != "" {
		db.Where("city = ?", params.City)
	}
	if params.Age != 0 {
		db.Where("age = ?", params.Age)
	}
	if params.Gender != 0 {
		db.Where("gender = ?", params.Gender)
	}
	if params.Phone != "" {
		db.Where("phone = ?", params.Phone)
	}
	if params.RealName != "" {
		db.Where("real_name = ?", params.RealName)
	}
	if params.Wechat != "" {
		db.Where("wechat = ?", params.Wechat)
	}

	err := db.Count(&total).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (d defaultCommunityActivitiesSignupModel) CancelSignup(signup *entity.CommunityActivitiesSignup) error {
	err := d.ctx.DB.Model(signup).Updates(&entity.CommunityActivitiesSignup{SignupStatus: SignupStatusFalse, CancelTime: time.Now()}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d defaultCommunityActivitiesSignupModel) FindOne(params FindOneActivitiesSignupParams) (*entity.CommunityActivitiesSignup, error) {
	var resp *entity.CommunityActivitiesSignup
	db := d.ctx.DB.Model(&entity.CommunityActivitiesSignup{})
	if params.Id != 0 {
		db.Where("id = ?", params.Id)
	}
	if params.TopicId != 0 {
		db.Where("topic_id = ?", params.TopicId)
	}
	if params.UserId != 0 {
		db.Where("user_id = ?", params.UserId)
	}
	if params.SignupStatus != 0 {
		db.Where("signup_status = ?", params.SignupStatus)
	}
	err := db.First(&resp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &entity.CommunityActivitiesSignup{}, nil
		}
		return &entity.CommunityActivitiesSignup{}, err
	}
	return resp, nil
}

func (d defaultCommunityActivitiesSignupModel) FindOneAPISignup(params FindOneActivitiesSignupParams) (*entity.APIActivitiesSignup, error) {
	var resp *entity.APIActivitiesSignup
	db := d.ctx.DB.Model(&entity.CommunityActivitiesSignup{})
	if params.Id != 0 {
		db.Where("id = ?", params.Id)
	}
	if params.TopicId != 0 {
		db.Where("topic_id = ?", params.TopicId)
	}
	if params.UserId != 0 {
		db.Where("user_id = ?", params.UserId)
	}

	err := db.Unscoped().Order("signup_time desc").First(&resp).Error
	if err != nil {
		return &entity.APIActivitiesSignup{}, err
	}
	return resp, nil
}

func (d defaultCommunityActivitiesSignupModel) FindAllAPISignup(params FindAllActivitiesSignupParams) ([]*entity.APIActivitiesSignup, int64, error) {
	list := make([]*entity.APIActivitiesSignup, 0)
	var total int64
	db := d.ctx.DB.Model(&entity.CommunityActivitiesSignup{}).
		Preload("User").
		Preload("Topic").
		Preload("Topic.User").
		Preload("Topic.Activity")

	if params.TopicId != 0 {
		db.Where("topic_id = ?", params.TopicId)
	}
	if params.UserId != 0 {
		db.Where("user_id = ?", params.UserId)
	}
	if params.City != "" {
		db.Where("city = ?", params.City)
	}
	if params.Age != 0 {
		db.Where("age = ?", params.Age)
	}
	if params.Gender != 0 {
		db.Where("gender = ?", params.Gender)
	}
	if params.Phone != "" {
		db.Where("phone = ?", params.Phone)
	}
	if params.RealName != "" {
		db.Where("real_name = ?", params.RealName)
	}
	if params.Wechat != "" {
		db.Where("wechat = ?", params.Wechat)
	}

	err := db.Count(&total).Offset(params.Offset).Limit(params.Limit).Order("signup_time desc").Unscoped().Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (d defaultCommunityActivitiesSignupModel) FindOneAPISignupById(id int64) (*entity.APIActivitiesSignup, error) {
	var resp *entity.APIActivitiesSignup
	err := d.ctx.DB.Model(&entity.CommunityActivitiesSignup{}).First(&resp, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &entity.APIActivitiesSignup{}, nil
		}
		return &entity.APIActivitiesSignup{}, err
	}
	return resp, nil
}

func (d defaultCommunityActivitiesSignupModel) Delete(id int64) error {
	if err := d.ctx.DB.Delete(&entity.CommunityActivitiesSignup{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (d defaultCommunityActivitiesSignupModel) Update(signup *entity.CommunityActivitiesSignup) error {
	err := d.ctx.DB.Save(signup).Error
	if err != nil {
		return err
	}
	return nil
}

func (d defaultCommunityActivitiesSignupModel) Create(signup *entity.CommunityActivitiesSignup) error {
	err := d.ctx.DB.Model(&entity.CommunityActivitiesSignup{}).
		WithContext(d.ctx.Context).
		Omit(clause.Associations).
		Omit("cancel_time").
		Create(signup).Error
	if err != nil {
		return err
	}
	return nil
}

func NewCommunityActivitiesSignupModel(ctx *mioContext.MioContext) ActivitiesSignupModel {
	return defaultCommunityActivitiesSignupModel{
		ctx: ctx,
	}
}
