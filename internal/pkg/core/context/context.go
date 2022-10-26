package context

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
)

type MioContext struct {
	//default context.Background()
	context.Context
	//default app.DB
	DB *gorm.DB
}

type Option func(mioctx *MioContext)

func WithDB(db *gorm.DB) Option {
	return func(ctx *MioContext) {
		ctx.DB = db
	}
}
func WithContext(ctx context.Context) Option {
	return func(mioctx *MioContext) {
		mioctx.Context = ctx
	}
}

func NewMioContext(options ...Option) *MioContext {
	mioContext := &MioContext{}
	for _, option := range options {
		option(mioContext)
	}
	if mioContext.Context == nil {
		mioContext.Context = context.Background()
	}
	if mioContext.DB == nil {
		mioContext.DB = app.DB
	}
	return mioContext
}

func NewBusinessContext(options ...Option) *MioContext {
	mioContext := &MioContext{}
	for _, option := range options {
		option(mioContext)
	}
	if mioContext.Context == nil {
		mioContext.Context = context.Background()
	}
	if mioContext.DB == nil {
		mioContext.DB = app.BusinessDB
	}
	return mioContext
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
