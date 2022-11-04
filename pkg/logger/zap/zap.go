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

// LogApplication 应用日志 框架自带
var LogApplication = zap.String("scene", "application_log")

// LogOperation 操作日志 自己打得
var LogOperation = zap.String("scene", "operation_log")

// LogAccess 访问日志
var LogAccess = zap.String("scene", "access_log")

// LogDatabase 数据库日志
var LogDatabase = zap.String("scene", "database")

var LevelMap = map[string]zapcore.Level{
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

func NewFileLogger(config LoggerConfig, opts ...zap.Option) *zap.Logger {
	logLevel := LevelMap[config.Level]
	encoder := getEncoder()
	writer := getWriter(config.Path, config.FileName, config.MaxSize)
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), zapcore.AddSync(writer), logLevel)
	opts = append(opts, zap.AddCaller())
	return zap.New(core, opts...)
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
func NewConsoleLogger(level string, opts ...zap.Option) *zap.Logger {
	if level == "" {
		level = "info"
	}

	levStr := strings.ToLower(level)
	var ok bool
	lev, ok := LevelMap[levStr]
	if !ok {
		lev = zapcore.InfoLevel
	}

	encoder := getEncoder()
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), zapcore.AddSync(os.Stdout), lev)
	opts = append(opts, zap.AddCaller())
	return zap.New(core, opts...)
}
