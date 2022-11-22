package aliyunzero

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"runtime"
	"time"
)
import "github.com/aliyun/aliyun-log-go-sdk/producer"

// LogTopicApplication 应用日志 框架自带
var LogTopicApplication = logx.Field("__topic__", "application_log")

// LogTopicOperation 操作日志 自己打得
var LogTopicOperation = logx.Field("__topic__", "operation_log")

// LogTopicAccess 访问日志
var LogTopicAccess = logx.Field("__topic__", "access_log")

func LogCaller(skip int) logx.LogField {
	_, file, line, _ := runtime.Caller(skip)
	return logx.Field("__caller", fmt.Sprintf("file:%s:%d", file, line))
}

type Option struct {
	Project  string
	LogStore string
}
type AlyWriter struct {
	producer *producer.Producer
	opts     Option
	fields   []logx.LogField
}

func NewAlyWriter(producer *producer.Producer, option Option) *AlyWriter {
	return &AlyWriter{producer: producer, opts: option}
}
func (l *AlyWriter) With(fields ...logx.LogField) *AlyWriter {
	w := l.clone()
	w.fields = append(w.fields, fields...)
	return w
}
func (l *AlyWriter) clone() *AlyWriter {
	return &AlyWriter{
		producer: l.producer,
		opts:     l.opts,
		fields:   l.fields,
	}
}
func (l AlyWriter) Alert(v interface{}) {

	l.Writer("alert", v)
}

func (l AlyWriter) Close() error {
	//l.producer.Close(5000)  //开始关闭producer的时候开始计时，超过传递的设定值还未能完全关闭producer的话会强制退出producer，此时可能会有部分数据未被成功发送而丢失
	l.producer.SafeClose() //producer提供了两种关闭模式，分为有限关闭和安全关闭，安全关闭会等待producer中缓存的所有的数据全部发送完成以后在关闭producer
	return nil
}
func (l AlyWriter) Writer(level string, v interface{}, fields ...logx.LogField) {
	logData := make(map[string]string)
	topic := ""
	fields = append(fields, l.fields...)
	for _, field := range fields {
		switch field {
		case LogTopicApplication, LogTopicAccess, LogTopicOperation:
			topic = field.Value.(string)
			continue
		}

		logData[field.Key] = fmt.Sprintf("%v", field.Value)
	}
	logData["__caller__"] = fmt.Sprintf("%v", LogCaller(6).Value)
	logData["__level__"] = level
	logData["__content__"] = fmt.Sprintf("%v", v)
	log := producer.GenerateLog(uint32(time.Now().Unix()), logData) //map[string]string{"content": "test", "content2": fmt.Sprintf("%v", "1111")}
	err := l.producer.SendLog(l.opts.Project, l.opts.LogStore, topic, os.Getenv("HOSTNAME"), log)
	if err != nil {
		fmt.Println("发送日志异常", err)
	}
}
func (l AlyWriter) Error(v interface{}, fields ...logx.LogField) {
	l.Writer("error", v, fields...)
}

func (l AlyWriter) Info(v interface{}, fields ...logx.LogField) {
	l.Writer("info", v, fields...)
}

//严峻的

func (l AlyWriter) Severe(v interface{}) {
	l.Writer("severe", v)
}

//缓慢的

func (l AlyWriter) Slow(v interface{}, fields ...logx.LogField) {
	l.Writer("slow", v, fields...)
}

//堆栈

func (l AlyWriter) Stack(v interface{}) {
	l.Writer("stack", v)
}

//统计

func (l AlyWriter) Stat(v interface{}, fields ...logx.LogField) {
	l.Writer("stat", v, fields...)
}
