package repository

import (
	"context"
	"gorm.io/gorm"
	mioctx "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type UserPlatformRepository struct {
	ctx *mioctx.MioContext
}

func (m *UserPlatformRepository) Delete(ctx context.Context, id int64) error {
	return m.ctx.DB.WithContext(ctx).Where("id = ?", id).Delete(&entity.UserPlatform{}).Error
}

func (m *UserPlatformRepository) FindOne(ctx context.Context, id int64) (*entity.UserPlatform, bool, error) {

	var resp entity.UserPlatform
	err := m.ctx.DB.Where("id = ?", id).Take(&resp).Error
	if err == nil {
		return &resp, true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return nil, false, err
}

func (m *UserPlatformRepository) FindOneByGuidPlatform(ctx context.Context, guid string, platform entity.UserPlatformType) (*entity.UserPlatform, bool, error) {

	var resp entity.UserPlatform
	err := m.ctx.DB.Model(&resp).Where("guid = ? and platform = ?", guid, platform).First(&resp).Error
	if err == nil {
		return &resp, true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return nil, false, err

}

func (m *UserPlatformRepository) FindOneByOpenidPlatform(ctx context.Context, openid string, platform entity.UserPlatform) (*entity.UserPlatform, bool, error) {

	var resp entity.UserPlatform
	err := m.ctx.DB.Model(&resp).Where("openid = ? and platform = ?", openid, platform).First(&resp).Error
	if err == nil {
		return &resp, true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return nil, false, err

}

func (m *UserPlatformRepository) Insert(ctx context.Context, data *entity.UserPlatform) error {

	return m.ctx.DB.WithContext(ctx).Create(data).Error

}

func (m *UserPlatformRepository) InsertBatch(ctx context.Context, data *[]entity.UserPlatform, batchSize int) error {
	err := m.ctx.DB.WithContext(ctx).CreateInBatches(data, batchSize).Error
	return err
}

func (m *UserPlatformRepository) Update(ctx context.Context, newData *entity.UserPlatform) error {

	return m.ctx.DB.WithContext(ctx).Save(newData).Error

}

// UpdateColumn 更新一列数据
// id 主键
// column 列名
// val 列值
// skipHook 是否跳过 Hook 方法且不追踪更新时间 true跳过 false不跳过 默认不跳过
func (m *UserPlatformRepository) UpdateColumn(ctx context.Context, id int64, column string, val interface{}, skipHook ...bool) error {
	var err error

	if len(skipHook) > 0 && skipHook[0] {
		err = m.ctx.DB.WithContext(ctx).Model(entity.UserPlatform{}).Where("id = ?", id).UpdateColumn(column, val).Error
	} else {
		err = m.ctx.DB.WithContext(ctx).Model(entity.UserPlatform{}).Where("id = ?", id).Update(column, val).Error
	}
	if err != nil {
		return err
	}

	return nil

}

// UpdateColumns 更新多列数据
// id 主键
// values map或者struct 当使用 struct 更新时，默认情况下，GORM 只会更新非零值的字段
// skipHook 是否跳过 Hook 方法且不追踪更新时间 true跳过 false不跳过 默认不跳过
func (m *UserPlatformRepository) UpdateColumns(ctx context.Context, id int64, values interface{}, skipHook ...bool) error {
	var err error

	if len(skipHook) > 0 && skipHook[0] {
		err = m.ctx.DB.WithContext(ctx).Model(entity.UserPlatform{}).Where("id = ?", id).UpdateColumns(values).Error
	} else {
		err = m.ctx.DB.WithContext(ctx).Model(entity.UserPlatform{}).Where("id = ?", id).Updates(values).Error
	}
	if err != nil {
		return err
	}

	return nil

}
