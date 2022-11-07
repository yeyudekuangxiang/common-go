package aliyunzero

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// FiledLog 可以初始化一个 logx.LogField 数组 每次记录日志都会把这些LogField一起记录上去
type FiledLog struct {
	logx.Logger
	fields []logx.LogField
}

func WithFiledLogger(logger logx.Logger, fields ...logx.LogField) *FiledLog {
	return &FiledLog{Logger: logger, fields: fields}
}

func (f FiledLog) Error(i ...interface{}) {
	f.Errorw(fmt.Sprint(i...))
}

func (f FiledLog) Errorf(s string, i ...interface{}) {
	f.Errorw(fmt.Errorf(s, i...).Error())
}

func (f FiledLog) Errorv(i interface{}) {
	f.Errorw(fmt.Sprintf("%v", i))
}

func (f FiledLog) Errorw(s string, field ...logx.LogField) {
	field = append(field, f.fields...)
	f.Logger.Errorw(s, field...)
}

func (f FiledLog) Info(i ...interface{}) {
	f.Infow(fmt.Sprint(i...))
}

func (f FiledLog) Infof(s string, i ...interface{}) {
	f.Infow(fmt.Sprintf(s, i...))
}

func (f FiledLog) Infov(i interface{}) {
	f.Infow(fmt.Sprintf("%v", i))
}

func (f FiledLog) Infow(s string, field ...logx.LogField) {
	field = append(field, f.fields...)
	f.Logger.Infow(s, field...)
}

func (f FiledLog) Slow(i ...interface{}) {
	f.Sloww(fmt.Sprint(i...))
}

func (f FiledLog) Slowf(s string, i ...interface{}) {
	f.Sloww(fmt.Sprintf(s, i...))
}

func (f FiledLog) Slowv(i interface{}) {
	f.Sloww(fmt.Sprintf("%v", i))
}

func (f FiledLog) Sloww(s string, field ...logx.LogField) {
	field = append(field, f.fields...)
	f.Logger.Sloww(s, field...)
}

func (f FiledLog) WithContext(ctx context.Context) logx.Logger {
	return f.Logger.WithContext(ctx)
}

func (f FiledLog) WithDuration(duration time.Duration) logx.Logger {
	return f.Logger.WithDuration(duration)
}
