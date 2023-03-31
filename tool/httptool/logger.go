package httptool

import (
	"log"
	"net/http"
	"time"
)

type Logger interface {
	Log(data LogData, err error)
}
type LogData struct {
	Url            string
	Header         http.Header
	ResponseHeader http.Header
	Start          time.Time
	//毫秒
	Duration     int64
	Method       string
	StatusCode   int
	RequestBody  []byte
	ResponseBody []byte
}

type ConsoleLogger struct {
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (d *ConsoleLogger) Log(data LogData, err error) {
	log.Printf(`data:%+v err:%+v\n`,
		data, err,
	)
}
