package {{.pkg}}

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
)

func New{{.appName}}Repository(db *gorm.DB, c cache.CacheConf) *{{.appName}}Repository {
	return &{{.appName}}Repository{
		newDefault{{.appName}}Repository(db, c),
	}
}

type {{.appName}}Repository struct {
	*default{{.appName}}Repository
}