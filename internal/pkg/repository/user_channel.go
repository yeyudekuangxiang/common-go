package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultUserChannelRepository = UserChannelRepository{DB: app.DB}

type UserChannelRepository struct {
	DB *gorm.DB
}

func (repo UserChannelRepository) FindByCid(by FindUserChannelBy) *entity.UserChannel {
	channel := entity.UserChannel{}
	db := repo.DB.Model(channel)
	if by.Cid != 0 {
		db.Where("cid = ?", by.Cid)
	}
	err := db.First(&channel).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return &channel
}

func (repo UserChannelRepository) FindByCode(by FindUserChannelBy) *entity.UserChannel {
	channel := entity.UserChannel{}
	db := repo.DB.Model(channel)
	if by.Cid != 0 {
		db.Where("cid != ?", by.Cid)
	}
	if by.Code != "" {
		db.Where("code = ?", by.Code)
	}
	err := db.First(&channel).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return &channel
}

func (repo UserChannelRepository) Save(channel *entity.UserChannel) error {
	return repo.DB.Save(channel).Error
}
func (repo UserChannelRepository) Create(channel *entity.UserChannel) error {
	return repo.DB.Create(channel).Error
}

func (repo UserChannelRepository) GetUserChannelPageList(by GetUserChannelPageListBy) (list []entity.UserChannel, total int64) {
	list = make([]entity.UserChannel, 0)
	db := repo.DB.Table("user_channel")
	if by.Pid > 0 {
		db.Where("pid = ?", by.Pid)
	}
	if by.Cid > 0 {
		db.Where("cid = ?", by.Cid)
	}
	if len(by.CidSlice) != 0 {
		db.Where("cid in (?)", by.CidSlice)
	}
	if by.Name != "" {
		db.Where("name like ?", "%"+by.Name+"%")
	}
	db.Select("*")
	db2 := repo.DB.Table("(?) as t", db)
	err := db2.Count(&total).
		Offset(by.Offset).
		Limit(by.Limit).
		Order("id desc").
		Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}
