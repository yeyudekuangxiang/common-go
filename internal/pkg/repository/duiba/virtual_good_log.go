package duiba

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/duiba"
)

var DefaultVirtualGoodLogRepository = VirtualGoodLogRepository{DB: app.DB}

type VirtualGoodLogRepository struct {
	DB *gorm.DB
}

func (repo VirtualGoodLogRepository) Save(log *duiba.VirtualGoodLog) error {
	return repo.DB.Save(log).Error
}
func (repo VirtualGoodLogRepository) Create(log *duiba.VirtualGoodLog) error {
	return repo.DB.Create(log).Error
}
func (repo VirtualGoodLogRepository) FindBy(by FindVirtualGoodLogBy) duiba.VirtualGoodLog {
	log := duiba.VirtualGoodLog{}
	db := repo.DB.Model(duiba.VirtualGoodLog{})

	if by.Params != "" {
		db.Where("params = ?", by.Params)
	}
	if by.OrderNum != "" {
		db.Where("order_num = ?", by.OrderNum)
	}

	err := db.Take(&log).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	return log
}
