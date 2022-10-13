package zap

import (
	"context"
	"github.com/medivhzhan/weapp/v3/logger"
	"go.uber.org/zap"
)

type WeappLogger struct {
	logger *zap.SugaredLogger
	level  logger.Level
}

func NewWeappLogger(logger *zap.SugaredLogger) *WeappLogger {
	return &WeappLogger{logger: logger}
}

func (w *WeappLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if w.level >= logger.Info {
		w.logger.Info(s, i)
	}
}

func (w *WeappLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if w.level >= logger.Warn {
		w.logger.Warn(s, i)
	}
}

func (w *WeappLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if w.level >= logger.Error {
		w.logger.Error(s, i)
	}
}

func (w *WeappLogger) SetLevel(level logger.Level) {
	w.level = level
}
