package repository

import (
	"database/sql"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultCommentRepository = NewCommentRepository(app.DB)

type (
	CommentModel interface {
		Insert(data *entity.CommentIndex) (*entity.CommentIndex, error)
		FindOne(id int64) (*entity.CommentIndex, error)
		FindOneQuery(builder *gorm.DB) (*entity.CommentIndex, error)
		FindCount(builder *gorm.DB) (int64, error)
		FindSum(builder *gorm.DB) (float64, error)
		FindAll(builder *gorm.DB, orderBy string) ([]*entity.CommentIndex, error)
		FindPageListByPage(builder *gorm.DB, offset, limit int64, orderBy string) ([]*entity.CommentIndex, error)
		FindPageListByIdDESC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CommentIndex, error)
		FindPageListByIdASC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CommentIndex, error)
		Delete(id, userId int64) error
		DeleteSoft(id, userId int64) error
		Update(data *entity.CommentIndex) error
		//UpdateWithVersion()
		Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
		RowBuilder() *gorm.DB
		CountBuilder(field string) *gorm.DB
		SumBuilder(field string) *gorm.DB
		AddTopicLikeCount(commentId int64, num int) error
	}

	defaultCommentRepository struct {
		Model *gorm.DB
	}
)

func NewCommentRepository(db *gorm.DB) CommentModel {
	return &defaultCommentRepository{
		Model: db,
	}
}

func (m *defaultCommentRepository) Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return m.Model.Model(&entity.CommentIndex{}).Transaction(fc, opts...)
}

func (m *defaultCommentRepository) RowBuilder() *gorm.DB {
	return m.Model.Model(&entity.CommentIndex{})
}

func (m *defaultCommentRepository) CountBuilder(field string) *gorm.DB {
	return m.Model.Model(&entity.CommentIndex{}).Select("COUNT(" + field + ")")
}

func (m *defaultCommentRepository) SumBuilder(field string) *gorm.DB {
	return m.Model.Model(&entity.CommentIndex{}).Select("SUM(" + field + ")")
}

func (m *defaultCommentRepository) Insert(data *entity.CommentIndex) (*entity.CommentIndex, error) {
	err := m.Model.Create(data).Error
	switch err {
	case nil:
		return data, nil
	default:
		return nil, err
	}
}

func (m *defaultCommentRepository) FindOne(id int64) (*entity.CommentIndex, error) {
	var resp entity.CommentIndex
	err := m.Model.First(&resp).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (m *defaultCommentRepository) FindOneQuery(builder *gorm.DB) (*entity.CommentIndex, error) {
	var resp entity.CommentIndex
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

func (m *defaultCommentRepository) FindCount(builder *gorm.DB) (int64, error) {
	var resp int64
	err := builder.Count(&resp).Error
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultCommentRepository) FindSum(builder *gorm.DB) (float64, error) {
	var resp float64
	err := builder.First(&resp).Error
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultCommentRepository) FindAll(builder *gorm.DB, orderBy string) ([]*entity.CommentIndex, error) {
	if orderBy == "" {
		builder.Order("comment_index.id DESC")
	} else {
		builder.Order(orderBy)
	}
	var resp []*entity.CommentIndex
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

func (m *defaultCommentRepository) FindPageListByPage(builder *gorm.DB, offset, limit int64, orderBy string) ([]*entity.CommentIndex, error) {
	if orderBy == "" {
		builder.Order("id DESC")
	} else {
		builder.Order(orderBy)
	}
	var resp []*entity.CommentIndex

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

func (m *defaultCommentRepository) FindPageListByIdDESC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CommentIndex, error) {
	if preMinId > 0 {
		builder = builder.Where(" id < ? ", preMinId)
	}
	var resp []*entity.CommentIndex
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

func (m *defaultCommentRepository) FindPageListByIdASC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CommentIndex, error) {
	if preMinId > 0 {
		builder = builder.Where(" id < ? ", preMinId)
	}
	var resp []*entity.CommentIndex
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

func (m *defaultCommentRepository) Delete(id, userId int64) error {
	var result entity.CommentIndex
	err := m.Model.Where("id = ? and member_id = ?", id, userId).First(&result).Error
	if err != nil {
		return err
	}
	err = m.Model.Transaction(func(tx *gorm.DB) error {
		if err = m.Model.Delete(result).Error; err != nil {
			return err
		}
		if err = m.Update(&result); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (m *defaultCommentRepository) DeleteSoft(id, userId int64) error {
	var result entity.CommentIndex
	err := m.Model.Where("id = ? and member_id = ?", id, userId).First(&result).Error
	if err != nil {
		return err
	}
	result.State = 1
	result.Count -= 1
	return m.Update(&result)
}

func (m *defaultCommentRepository) Update(data *entity.CommentIndex) error {
	var result entity.CommentIndex
	err := m.Model.Where("id = ? and member_id = ?", data.ID, data.MemberId).First(&result).Error
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
	return m.Model.Model(&result).Updates(&result).Error
}

func (m *defaultCommentRepository) AddTopicLikeCount(commentId int64, num int) error {
	db := m.Model.Model(&entity.CommentIndex{}).Where("id = ?", commentId)
	//避免点赞数为负数
	if num < 0 {
		db.Where("like_count >= ?", -num)
	}
	return db.Update("like_count", gorm.Expr("like_count + ?", num)).Error
}
