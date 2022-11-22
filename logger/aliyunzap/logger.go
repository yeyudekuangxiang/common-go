package aliyunzap

import (
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	DefaultEncoderConfig = EncoderConfig{
		MessageKey:  "__content__",
		TimeKey:     "__timestamp__",
		LevelKey:    "__level__",
		FunctionKey: "__function__",
		CallerKey:   "__caller__",
		EncodeTime: func(t time.Time) string {
			return t.Format(time.RFC3339Nano)
		},
		EncodeCaller: func(caller zapcore.EntryCaller) string {
			return caller.FullPath()
		},
	}
	DefaultEncoder = NewAliYunEncoder(&DefaultEncoderConfig)
)
