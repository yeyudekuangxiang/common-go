package {{.pkg}}

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
)

type Default{{.appName}}Repository struct {
	db               *gorm.DB
	c                cache.CacheConf
	{{.models}}
}

func (repo *Default{{.appName}}Repository) Transaction(f func(repoTx *Default{{.appName}}Repository) error) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		newRepo := newDefault{{.appName}}Repository(tx, repo.c)
		return f(newRepo)
	})
}

func newDefault{{.appName}}Repository(db *gorm.DB, c cache.CacheConf) *Default{{.appName}}Repository {
	return &Default{{.appName}}Repository{
		db:               db,
		c:                c,
		{{.newModels}}
	}
}
