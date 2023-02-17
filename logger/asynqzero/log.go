package asynqzero

import "github.com/zeromicro/go-zero/core/logx"

type Logger struct {
	log logx.Logger
}

func NewLogger(log logx.Logger) *Logger {
	return &Logger{log: log}
}

func (l Logger) Debug(args ...interface{}) {
	l.log.Info(args...)
}

func (l Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l Logger) Warn(args ...interface{}) {
	l.log.Error(args...)
}

func (l Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l Logger) Fatal(args ...interface{}) {
	l.log.Error(args...)
}
