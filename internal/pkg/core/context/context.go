package context

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type MioContext struct {
	context.Context
	DB *gorm.DB
}

// BeginTran begins a transaction
func (ctx *MioContext) BeginTran(options ...*sql.TxOptions) *MioContext {
	db := ctx.DB.Begin(options...)
	return &MioContext{
		Context: ctx.Context,
		DB:      db,
	}
}

// CommitTran commit a transaction
func (ctx *MioContext) CommitTran() {
	ctx.DB.Commit()
}

// RollbackTran rollback a transaction
func (ctx *MioContext) RollbackTran() {
	ctx.DB.Rollback()
}

// Transaction start a transaction as a block, return error will rollback, otherwise to commit.
func (ctx *MioContext) Transaction(f func(ctx *MioContext) error) error {
	return ctx.DB.Transaction(func(tx *gorm.DB) error {
		return f(&MioContext{
			Context: ctx.Context,
			DB:      tx,
		})
	})
}
