package {{.pkg}}

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Default{{.appName}}Repository struct {
	db               *gorm.DB
	c                cache.CacheConf
	{{.models}}
}


// Transaction 开启一个事务
func (repo *Default{{.appName}}Repository) Transaction(f func(repoTx *Default{{.appName}}Repository) error) error {
	return repo.db.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		newRepo := newDefault{{.appName}}Repository(tx, repo.c,  ModelOptionSkipCache())
		return f(newRepo)
	})
}

func newDefault{{.appName}}Repository(db *gorm.DB, c cache.CacheConf, opts ...modelOption) *Default{{.appName}}Repository {
	return &Default{{.appName}}Repository{
		db:               db,
		c:                c,
		{{.newModels}}
	}
}
