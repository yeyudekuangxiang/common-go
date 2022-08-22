package repository

import (
	"database/sql"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultCouponHistoryRepository = NewCouponHistoryRepository(app.DB)

type (
	CouponHistoryModel interface {
		Insert(data *entity.CouponHistory) (*entity.CouponHistory, error)
		FindOne(id int64) (*entity.CouponHistory, error)
		FindOneQuery(builder *gorm.DB) (*entity.CouponHistory, error)
		FindCount(builder *gorm.DB) (int64, error)
		FindSum(builder *gorm.DB) (float64, error)
		FindAll(builder *gorm.DB, orderBy string) ([]*entity.CouponHistory, error)
		FindPageListByPage(builder *gorm.DB, offset, limit int64, orderBy string) ([]*entity.CouponHistory, error)
		FindPageListByIdDESC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CouponHistory, error)
		FindPageListByIdASC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CouponHistory, error)
		//Delete(id, userId int64) error
		//DeleteSoft(id, userId int64) error
		Update(data *entity.CouponHistory) error
		//UpdateWithVersion()
		Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
		RowBuilder() *gorm.DB
		CountBuilder(field string) *gorm.DB
		SumBuilder(field string) *gorm.DB
	}

	defaultCouponHistoryRepository struct {
		Model *gorm.DB
	}
)

func (m *defaultCouponHistoryRepository) Update(data *entity.CouponHistory) error {
	var result entity.CommentIndex
	err := m.Model.Where("open_id = ?", data.OpenId).First(&result).Error
	if err != nil {
		return err
	}
	if data.Code != "" {
		result.Message = data.Code
	}
	return m.Model.Model(&result).Updates(&result).Error
}

func NewCouponHistoryRepository(db *gorm.DB) CouponHistoryModel {
	return &defaultCouponHistoryRepository{
		Model: db,
	}
}

func (m *defaultCouponHistoryRepository) Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return m.Model.Model(&entity.CouponHistory{}).Transaction(fc, opts...)
}

func (m *defaultCouponHistoryRepository) RowBuilder() *gorm.DB {
	return m.Model.Model(&entity.CouponHistory{})
}

func (m *defaultCouponHistoryRepository) CountBuilder(field string) *gorm.DB {
	return m.Model.Model(&entity.CouponHistory{}).Select("COUNT(" + field + ")")
}

func (m *defaultCouponHistoryRepository) SumBuilder(field string) *gorm.DB {
	return m.Model.Model(&entity.CouponHistory{}).Select("SUM(" + field + ")")
}

func (m *defaultCouponHistoryRepository) Insert(data *entity.CouponHistory) (*entity.CouponHistory, error) {
	err := m.Model.Create(data).Error
	switch err {
	case nil:
		return data, nil
	default:
		return nil, err
	}
}

func (m *defaultCouponHistoryRepository) FindOne(id int64) (*entity.CouponHistory, error) {
	var resp entity.CouponHistory
	err := m.Model.First(&resp, id).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (m *defaultCouponHistoryRepository) FindOneQuery(builder *gorm.DB) (*entity.CouponHistory, error) {
	var resp entity.CouponHistory
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

func (m *defaultCouponHistoryRepository) FindCount(builder *gorm.DB) (int64, error) {
	var resp int64
	err := builder.Find(&resp).Error
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultCouponHistoryRepository) FindSum(builder *gorm.DB) (float64, error) {
	var resp float64
	err := builder.First(&resp).Error
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultCouponHistoryRepository) FindAll(builder *gorm.DB, orderBy string) ([]*entity.CouponHistory, error) {
	if orderBy == "" {
		builder.Order("comment_index.id DESC")
	} else {
		builder.Order(orderBy)
	}
	var resp []*entity.CouponHistory
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

func (m *defaultCouponHistoryRepository) FindPageListByPage(builder *gorm.DB, offset, limit int64, orderBy string) ([]*entity.CouponHistory, error) {
	if orderBy == "" {
		builder.Order("id DESC")
	} else {
		builder.Order(orderBy)
	}
	var resp []*entity.CouponHistory

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

func (m *defaultCouponHistoryRepository) FindPageListByIdDESC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CouponHistory, error) {
	if preMinId > 0 {
		builder = builder.Where(" id < ? ", preMinId)
	}
	var resp []*entity.CouponHistory
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

func (m *defaultCouponHistoryRepository) FindPageListByIdASC(builder *gorm.DB, preMinId, limit int64) ([]*entity.CouponHistory, error) {
	if preMinId > 0 {
		builder = builder.Where(" id < ? ", preMinId)
	}
	var resp []*entity.CouponHistory
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
