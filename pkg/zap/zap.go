package zap

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"mio/pkg/wxwork"
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
	return zap.New(core, zap.AddCaller(), zap.Hooks(func(entry zapcore.Entry) error {
		if entry.Level >= zapcore.ErrorLevel {
			err := wxwork.SendRobotMessage("f0edb1a2-3f9b-4a5d-aa15-9596a32840ec", wxwork.Markdown{
				Content: fmt.Sprintf(
					"**来源:**日志 \n\n**level:**%s \n\n**time**:%s \n\n**message**:%s \n\n**caller**:%+v \n\n**stack**:%s", entry.Level, entry.Time.Format("2006-01-02 15:04:05"), entry.Message, entry.Caller, entry.Stack),
			})
			if err != nil {
				log.Printf("推送异常到企业微信失败 %+v %v", entry, err)
			}
		}
		return nil
	}))
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
