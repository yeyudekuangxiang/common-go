package {{.pkg}}

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
)

type default{{.appName}}Repository struct {
	db               *gorm.DB
	c                cache.CacheConf
	{{.models}}
}

func (repo *default{{.appName}}Repository) Transaction(f func(repoTx *default{{.appName}}Repository) error) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		newRepo := newDefault{{.appName}}Repository(tx, repo.c)
		return f(newRepo)
	})
}

func newDefault{{.appName}}Repository(db *gorm.DB, c cache.CacheConf) *default{{.appName}}Repository {
	return &default{{.appName}}Repository{
		db:               db,
		c:                c,
		{{.newModels}}
	}
}
