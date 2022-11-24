package repository

import (
	"database/sql"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type (
	CarbonCommentModel interface {
		Insert(data *entity.CarbonCommentIndex) (*entity.CarbonCommentIndex, error)
		FindOne(id int64) (*entity.CarbonCommentIndex, error)
		FindOneQuery(builder *gorm.DB) (*entity.CarbonCommentIndex, error)
		FindCount(builder *gorm.DB) (int64, error)
		FindSum(builder *gorm.DB) (float64, error)
		FindAll(builder *gorm.DB, orderBy string) ([]*entity.CarbonCommentIndex, error)
		FindPageListByPage(builder *gorm.DB, offset, limit int64, orderBy string) ([]*entity.CarbonCommentIndex, error)
		FindPageListByIdDESC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CarbonCommentIndex, error)
		FindPageListByIdASC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CarbonCommentIndex, error)
		Delete(id, userId int64) error
		DeleteSoft(id, userId int64) error
		Update(data *entity.CarbonCommentIndex) error
		//UpdateWithVersion()
		Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
		RowBuilder() *gorm.DB
		CountBuilder(field string) *gorm.DB
		SumBuilder(field string) *gorm.DB
		AddCommentLikeCount(commentId int64, num int) error
		FindListByIds(ids []int64) []*entity.CarbonCommentIndex
	}

	defaultCarbonCommentModel struct {
		ctx *mioContext.MioContext
	}
)

func (m *defaultCarbonCommentModel) FindListByIds(ids []int64) []*entity.CarbonCommentIndex {
	commentList := make([]*entity.CarbonCommentIndex, len(ids))
	err := app.DB.Model(&entity.CarbonCommentIndex{}).
		Where("id in (?)", ids).
		//Where("state = ?", 0).
		Find(&commentList).Error
	if err != nil {
		return []*entity.CarbonCommentIndex{}
	}
	return commentList
}

func NewCarbonCommentModel(ctx *mioContext.MioContext) CarbonCommentModel {
	return &defaultCarbonCommentModel{
		ctx: ctx,
	}
}

func (m *defaultCarbonCommentModel) Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return m.ctx.DB.Transaction(fc, opts...)
}

func (m *defaultCarbonCommentModel) RowBuilder() *gorm.DB {
	return m.ctx.DB.WithContext(m.ctx.Context).Model(&entity.CarbonCommentIndex{})
}

func (m *defaultCarbonCommentModel) CountBuilder(field string) *gorm.DB {
	return m.ctx.DB.Model(&entity.CarbonCommentIndex{}).Select("COUNT(" + field + ")")
}

func (m *defaultCarbonCommentModel) SumBuilder(field string) *gorm.DB {
	return m.ctx.DB.Model(&entity.CarbonCommentIndex{}).Select("SUM(" + field + ")")
}

func (m *defaultCarbonCommentModel) Insert(data *entity.CarbonCommentIndex) (*entity.CarbonCommentIndex, error) {
	err := m.ctx.DB.Create(data).Error
	switch err {
	case nil:
		return data, nil
	default:
		return nil, err
	}
}

func (m *defaultCarbonCommentModel) FindOne(id int64) (*entity.CarbonCommentIndex, error) {
	var resp entity.CarbonCommentIndex
	err := m.ctx.DB.First(&resp, id).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (m *defaultCarbonCommentModel) FindOneQuery(builder *gorm.DB) (*entity.CarbonCommentIndex, error) {
	var resp entity.CarbonCommentIndex
	err := builder.First(&resp).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (m *defaultCarbonCommentModel) FindCount(builder *gorm.DB) (int64, error) {
	var resp int64
	err := builder.Count(&resp).Error
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultCarbonCommentModel) FindSum(builder *gorm.DB) (float64, error) {
	var resp float64
	err := builder.First(&resp).Error
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultCarbonCommentModel) FindAll(builder *gorm.DB, orderBy string) ([]*entity.CarbonCommentIndex, error) {
	if orderBy == "" {
		builder.Order("comment_index.id DESC")
	} else {
		builder.Order(orderBy)
	}
	var resp []*entity.CarbonCommentIndex
	err := builder.Find(&resp).Error
	switch err {
	case nil:
		return resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (m *defaultCarbonCommentModel) FindPageListByPage(builder *gorm.DB, offset, limit int64, orderBy string) ([]*entity.CarbonCommentIndex, error) {
	if orderBy == "" {
		builder.Order("id DESC")
	} else {
		builder.Order(orderBy)
	}
	var resp []*entity.CarbonCommentIndex

	err := builder.Offset(int(offset)).Limit(int(limit)).Find(&resp).Error
	switch err {
	case nil:
		return resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (m *defaultCarbonCommentModel) FindPageListByIdDESC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CarbonCommentIndex, error) {
	if preMinId > 0 {
		builder = builder.Where(" id < ? ", preMinId)
	}
	var resp []*entity.CarbonCommentIndex
	err := builder.Order("id DESC").Limit(int(limit)).Find(&resp).Error
	switch err {
	case nil:
		return resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (m *defaultCarbonCommentModel) FindPageListByIdASC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CarbonCommentIndex, error) {
	if preMinId > 0 {
		builder = builder.Where(" id < ? ", preMinId)
	}
	var resp []*entity.CarbonCommentIndex
	err := builder.Order("id ASC").Limit(int(limit)).Find(&resp).Error
	switch err {
	case nil:
		return resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (m *defaultCarbonCommentModel) Delete(id, userId int64) error {
	var result entity.CarbonCommentIndex
	err := m.ctx.DB.Where("id = ? and member_id = ?", id, userId).First(&result).Error
	if err != nil {
		return err
	}
	if err = m.ctx.DB.Delete(result).Error; err != nil {
		return err
	}
	return nil
}

func (m *defaultCarbonCommentModel) DeleteSoft(id, userId int64) error {
	var result entity.CarbonCommentIndex
	err := m.ctx.DB.Where("id = ? and member_id = ?", id, userId).First(&result).Error
	if err != nil {
		return err
	}
	result.State = 1
	result.Count -= 1
	return m.Update(&result)
}

func (m *defaultCarbonCommentModel) Update(data *entity.CarbonCommentIndex) error {
	var result entity.CarbonCommentIndex
	err := m.ctx.DB.Where("id = ? and member_id = ?", data.Id, data.MemberId).First(&result).Error
	if err != nil {
		return err
	}
	if data.Message != "" {
		result.Message = data.Message
	}
	if data.Floor != 0 {
		result.Floor = data.Floor
	}
	if data.RootCount != 0 {
		result.RootCount = data.RootCount
	}
	if data.Count != 0 {
		result.Count = data.Count
	}
	if data.Attrs >= 0 {
		result.Attrs = data.Attrs
	}
	if data.State >= 0 {
		result.State = data.State
	}
	if data.LikeCount != 0 {
		result.LikeCount = data.LikeCount
	}
	if data.HateCount != 0 {
		result.HateCount = data.HateCount
	}
	return m.ctx.DB.Model(&result).Updates(&result).Error
}

func (m *defaultCarbonCommentModel) AddCommentLikeCount(commentId int64, num int) error {
	db := m.ctx.DB.Model(&entity.CarbonCommentIndex{}).Where("id = ?", commentId)
	//避免点赞数为负数
	if num < 0 {
		db.Where("like_count >= ?", -num)
	}
	return db.Update("like_count", gorm.Expr("like_count + ?", num)).Error
}
