package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.miotech.com/miotech-application/backend/common-go/logger/aliyunzap"
	"gitlab.miotech.com/miotech-application/backend/common-go/logger/aliyunzero"
	mzap "gitlab.miotech.com/miotech-application/backend/common-go/logger/zap"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/aliyuntool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"gitlab.miotech.com/miotech-application/backend/common-go/wxwork"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/util"
	"os"
	"time"
)

func InitLog() {
	logx.DisableStat()
	cred, err := aliyuntool.UpdateAliYunStsTokenFunc(nil)
	if err != nil && err != aliyuntool.ErrCredentialNotFound {
		log.Panic("初始化阿里云UpdateAliYunStsTokenFunc异常", err)
	}

	if (config.Config.AliLog.AccessKey != "" && config.Config.AliLog.AccessSecret != "") || cred != nil {
		log.Println("发现阿里云日志配置,自动初始化阿里云日志")
		InitAliyunLog(cred)
	} else {
		log.Println("未发现阿里云日志配置,自动初始化控制台日志")
		InitConsoleLog()
	}
}
func InitConsoleLog() {
	log.Println("初始化控制台日志组件...")
	logger := mzap.NewConsoleLogger("debug", zap.AddCaller(), wxRobotHook)
	*app.OriginLogger = *logger
	s := logger.WithOptions(zap.IncreaseLevel(mzap.LevelMap[config.Config.Log.Level])).With(mzap.LogOperation).Sugar()
	*app.Logger = *s
	log.Println("初始化控制台日志组件成功")
}
func InitFileLog() {
	log.Println("初始化文件日志组件...")
	var loggerConfig mzap.LoggerConfig
	var err error
	err = util.MapTo(config.Config.Log, &loggerConfig)
	if err != nil {
		log.Panic(err)
	}
	loggerConfig.Path = "runtime"
	loggerConfig.FileName = "log.log"
	logger := mzap.NewFileLogger(mzap.LoggerConfig{
		Level:    "debug",
		Path:     loggerConfig.Path,
		FileName: loggerConfig.FileName,
		MaxSize:  loggerConfig.MaxSize,
	}, zap.AddCaller(), wxRobotHook)
	*app.OriginLogger = *logger

	s := logger.WithOptions(zap.IncreaseLevel(mzap.LevelMap[loggerConfig.Level])).With(mzap.LogOperation).Sugar()
	*app.Logger = *s
	log.Println("初始化文件日志组件成功")
}
func InitAliyunLog(UpdateStsToken func() (accessKeyID, accessKeySecret, securityToken string, expireTime time.Time, err error)) {
	log.Println("初始化阿里云日志组件...")

	aliConfig := config.Config.AliLog
	prd := producer.InitProducer(&producer.ProducerConfig{
		Endpoint:         aliConfig.Endpoint,
		AccessKeyID:      aliConfig.AccessKey,
		AccessKeySecret:  aliConfig.AccessSecret,
		UpdateStsToken:   UpdateStsToken,
		StsTokenShutDown: make(chan struct{}),
	})
	prd.Start()

	//替换zero log
	aliWriter := aliyunzero.NewAlyWriter(prd, aliyunzero.Option{
		Project:  aliConfig.ProjectName,
		LogStore: aliConfig.LogStore,
	}).With(aliyunzero.LogTopicOperation)
	logx.SetWriter(aliWriter)

	//zap logger
	core := aliyunzap.NewAliYunCore(aliyunzap.DefaultEncoder, prd, aliyunzap.LogConfig{
		ProjectName:  aliConfig.ProjectName,
		LogStore:     aliConfig.LogStore,
		Topic:        aliConfig.Topic,
		Source:       aliConfig.Source,
		LevelEnabler: zap.DebugLevel,
	})
	logger := zap.New(core, zap.AddCaller(), wxRobotHook)
	*app.OriginLogger = *logger

	s := logger.WithOptions(zap.IncreaseLevel(mzap.LevelMap[config.Config.Log.Level])).With(mzap.LogOperation).Sugar()
	*app.Logger = *s
	log.Println("初始化阿里云日志组件成功")
}
func closeLogger() {
	if app.Logger != nil {
		log.Println("关闭日志")
		err := app.Logger.Sync()
		if err != nil {
			log.Println("关闭日志失败", err)
		} else {
			log.Println("关闭日志成功")
		}
	}

}
func InitHttpToolLog() {
	httptool.SetLogger(httptool.FuncLogger(func(data httptool.LogData, err error) {
		dataByte, err2 := json.Marshal(data)
		if err2 != nil {
			logx.Error("转化http数据异常 data:%+v err:%+v err2:%+v", data, err, err2)
		} else {
			logx.Infof(`request_log data:%+v err:%+v`, string(dataByte), err)
		}
	}))
}

var wxRobotHook = zap.Hooks(func(entry zapcore.Entry) error {
	if entry.Level >= zapcore.ErrorLevel {
		if config.Config.App.Env != "prod" {
			return nil
		}
		err := wxwork.SendRobotMessage(config.Constants.WxWorkBugRobotKey, wxwork.Markdown{
			Content: fmt.Sprintf(
				"**容器:**%s \n\n**来源:**日志 \n\n**level:**%s \n\n**time**:%s \n\n**message**:%s \n\n**caller**:%+v \n\n**stack**:%s", os.Getenv("HOSTNAME"), entry.Level, entry.Time.Format("2006-01-02 15:04:05"), entry.Message, entry.Caller, entry.Stack),
		})
		if err != nil {
			log.Printf("推送异常到企业微信失败 %+v %v", entry, err)
		}
	}
	return nil
})
