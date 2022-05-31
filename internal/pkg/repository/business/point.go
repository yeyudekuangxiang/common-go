package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
)

var DefaultPointRepository = PointRepository{DB: app.DB}

type PointRepository struct {
	DB *gorm.DB
}
