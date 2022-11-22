package aliyunzap

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/gogo/protobuf/proto"
	"gitlab.miotech.com/miotech-application/backend/common-go/pkg/logger/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type LogConfig struct {
	ProjectName string
	LogStore    string
	Topic       string
	Source      string
	zapcore.LevelEnabler
}

func NewAliYunCore(enc *AliYunEncoder, producer *producer.Producer, config LogConfig) *AliYunCore {
	return &AliYunCore{
		enc:       enc,
		producer:  producer,
		LogConfig: config,
	}
}

type AliYunCore struct {
	LogConfig
	enc      *AliYunEncoder
	producer *producer.Producer
	client   sls.ClientInterface
}

func (a *AliYunCore) With(fields []zapcore.Field) zapcore.Core {

	clone := a.clone()
	for _, f := range fields {
		switch f {
		case zap.LogAccess, zap.LogApplication, zap.LogDatabase, zap.LogOperation:
			clone.Topic = f.String
			continue
		}
		addFields(clone.enc, []zapcore.Field{f})
	}

	return clone
}
func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
func (a *AliYunCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if a.Enabled(ent.Level) {
		return ce.AddCore(ent, a)
	}
	return ce
}

func (a *AliYunCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	m, err := a.enc.EncodeEntry(entry, fields)
	if err != nil {
		return err
	}
	var content []*sls.LogContent

	for key, value := range m {
		content = append(content, &sls.LogContent{
			Key:   proto.String(key),
			Value: proto.String(value),
		})
	}

	log := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: content,
	}

	source := a.Source
	if source == "" {
		source = os.Getenv("HOSTNAME")
	}

	return a.producer.SendLog(a.ProjectName, a.LogStore, a.Topic, source, log)

}

func (a *AliYunCore) Sync() error {
	return a.producer.Close(30)
}
func (a *AliYunCore) clone() *AliYunCore {
	return &AliYunCore{
		LogConfig: a.LogConfig,
		client:    a.client,
		producer:  a.producer,
		enc:       a.enc.clone(),
	}
}
