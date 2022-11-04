package aliyunzero

import (
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"reflect"
)

type AliYunSlsConf struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Project         string
	LogStore        string
}

func AutoInitAliyunLog(conf interface{}) {
	val := reflect.ValueOf(conf)
	v := val.FieldByName("AliYunSlsConf")
	if !v.IsValid() {
		logx.Alert("未启用阿里云日志 缺少AliYunSlsConf配置项")
		return
	}

	// 文档 https://help.aliyun.com/document_detail/48874.html
	logConf, ok := v.Interface().(AliYunSlsConf)
	if !ok {
		logx.Alert("未启用阿里云日志 配置AliYunSlsConf类型不是aliyunzero.AliYunSlsConf")
		return
	}
	if logConf.AccessKeyID == "" || logConf.AccessKeySecret == "" {
		logx.Alert("未启用阿里云日志 配置缺少AccessKeyID或者AccessKeySecret")
		return
	}

	pdr := producer.InitProducer(&producer.ProducerConfig{
		Endpoint:        logConf.Endpoint,
		AccessKeyID:     logConf.AccessKeyID,
		AccessKeySecret: logConf.AccessKeySecret,
		//MaxBatchCount:   1,
	})
	pdr.Start()
	ali := NewAlyWriter(pdr, Option{
		Project:  logConf.Project,
		LogStore: logConf.LogStore,
	})
	logx.SetWriter(ali)
	log.Println("成功将日志输出替换成阿里云日志")
	return
}
