package entity

import "gorm.io/gorm"

type OrderByList []OrderBy
type OrderBy string

var ErrNotFount = gorm.ErrRecordNotFound
