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

// Policy 设置从主库还是从库读
func (repo *Default{{.appName}}Repository) Policy(operation dbresolver.Operation) *Default{{.appName}}Repository  {
	db := repo.db.Clauses(operation).Session(&gorm.Session{NewDB: true})
	return newDefault{{.appName}}Repository(db, repo.c)
}

// Transaction 开启一个事务
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
