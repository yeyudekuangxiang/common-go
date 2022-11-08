package zap

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type GormLogger struct {
	level         logger.LogLevel
	logger        *zap.SugaredLogger
	SlowThreshold time.Duration
}

func NewGormLogger(sloger *zap.SugaredLogger) *GormLogger {
	return &GormLogger{logger: sloger, SlowThreshold: time.Second * 3, level: logger.Info}
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	g.level = level
	return g
}

func (g *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if g.level >= logger.Info {
		g.logger.Infof(infoStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...)
	}
}

func (g *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if g.level >= logger.Warn {
		g.logger.Warnf(warnStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...)
	}
}

func (g *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if g.level >= logger.Error {
		g.logger.Errorf(errStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...)
	}
}

var (
	infoStr      = "%s\n[info] "
	warnStr      = "%s\n[warn] "
	errStr       = "%s\n[error] "
	traceStr     = "%s\n[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
)

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && g.level >= logger.Error:
		sql, rows := fc()
		if rows == -1 {
			g.logger.Debugf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.logger.Debugf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > g.SlowThreshold && g.level >= logger.Warn && g.SlowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
		if rows == -1 {
			g.logger.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.logger.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case g.level == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			g.logger.Infof(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.logger.Infof(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
