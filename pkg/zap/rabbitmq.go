package zap

import "go.uber.org/zap"

type RabbitmqLogger struct {
	*zap.SugaredLogger
}

func NewRabbitmqLogger(sugaredLogger *zap.SugaredLogger) *RabbitmqLogger {
	return &RabbitmqLogger{SugaredLogger: sugaredLogger}
}

func (r RabbitmqLogger) Tracef(s string, i ...interface{}) {
	r.SugaredLogger.Debugf(s, i)
}
