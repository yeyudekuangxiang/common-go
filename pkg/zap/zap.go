package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

type LoggerConfig struct {
	Level    string //debug info warn error
	Path     string
	FileName string
	MaxSize  int
}

func NewZapLogger(config LoggerConfig) *zap.Logger {
	logLevel := levelMap[config.Level]
	encoder := getEncoder()
	writer := getWriter(config.Path, config.FileName, config.MaxSize)
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), zapcore.AddSync(writer), logLevel)
	return zap.New(core, zap.AddCaller())
}
func getEncoder() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
func getWriter(outputPath, fileName string, maxSize int) io.Writer {
	filename := filepath.Join(outputPath, fileName)
	outputPath = outputPath + string(os.PathSeparator)
	return &lumberjack.Logger{
		Filename:  filename,
		MaxSize:   maxSize, // megabytes,
		LocalTime: true,
	}
}
func DefaultLogger(level ...string) *zap.Logger {
	lev := zapcore.InfoLevel
	if len(level) > 0 {
		levStr := strings.ToLower(level[0])
		var ok bool
		lev, ok = levelMap[levStr]
		if !ok {
			lev = zapcore.InfoLevel
		}
	}

	encoder := getEncoder()
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), zapcore.AddSync(os.Stdout), lev)
	return zap.New(core, zap.AddCaller())
}
