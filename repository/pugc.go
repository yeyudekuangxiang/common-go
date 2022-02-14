package repository

import (
	"gorm.io/gorm"
	"mio/core/app"
	Pugc "mio/model/pugc"
)

var DefaultPugcRepository IPugcRepository = NewPugcRepository()

type IPugcRepository interface {
	// Insert GetPugcById 根据用id获取用户信息
	Insert(pugc *Pugc.PugcAddModel) error
}

func NewPugcRepository() PugcRepository {
	return PugcRepository{}
}

type PugcRepository struct {
}

func (u PugcRepository) Insert(pugc *Pugc.PugcAddModel) error {
	if err := app.DB.Table("pugc").Create(&pugc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
