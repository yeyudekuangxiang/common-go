package config

var Config = app{
	App:              appSetting{},
	Http:             httpSetting{},
	Database:         databaseSetting{},
	Log:              logSetting{},
	AliLog:           aliLogSetting{},
	Weapp:            wxSetting{},
	MioSubOA:         wxSetting{},
	MioSrvOA:         wxSetting{},
	Redis:            redisSetting{},
	DuiBa:            duiBaSetting{},
	OSS:              ossSetting{},
	BaiDu:            baiDuSetting{},
	Java:             javaConfig{},
	Zhuge:            zhugeConfig{},
	AMQP:             amqpSetting{},
	CouponRpc:        rpcSetting{},
	ActivityZyh:      activityZyh{},
	BaiDuImageSearch: baiDuImageSearchSetting{},
	TokenCenterRpc:   rpcSetting{},
	ActivityRpc:      rpcSetting{},
	DatabaseBusiness: databaseSetting{},
	Sms:              sms{},
	SmsMarket:        smsMarket{},
	Prometheus:       promSetting{},
}

type app struct {
	App              appSetting              `ini:"app"`
	Http             httpSetting             `ini:"http"`
	Database         databaseSetting         `ini:"database"`
	DatabaseBusiness databaseSetting         `ini:"databaseBusiness"`
	Log              logSetting              `ini:"log"`
	AliLog           aliLogSetting           `ini:"aliLog"`
	Weapp            wxSetting               `ini:"weapp"`
	MioSubOA         wxSetting               `ini:"mioSubOa"` //绿喵订阅号配置
	MioSrvOA         wxSetting               `ini:"mioSrvOa"` //绿喵服务号配置
	Redis            redisSetting            `ini:"redis"`
	DuiBa            duiBaSetting            `ini:"duiba"`
	OSS              ossSetting              `ini:"oss"`
	AMQP             amqpSetting             `ini:"amqp"`
	BaiDu            baiDuSetting            `ini:"baidu"`
	Java             javaConfig              `ini:"java"`
	Zhuge            zhugeConfig             `ini:"zhuge"`
	CouponRpc        rpcSetting              `ini:"couponRpc"`
	ActivityZyh      activityZyh             `ini:"activityZyh"`
	BaiDuImageSearch baiDuImageSearchSetting `ini:"baiduImageSearch"`
	TokenCenterRpc   rpcSetting              `ini:"tokenCenterRpc"`
	ActivityRpc      rpcSetting              `ini:"activityRpc"`
	PointRpc         rpcSetting              `ini:"pointRpc"`
	Sms              sms                     `ini:"sms"`
	SmsMarket        smsMarket               `ini:"smsMarket"`
	Prometheus       promSetting             `ini:"prometheus"`
}

type appSetting struct {
	TokenKey string
	Domain   string
	Debug    bool
	//prod dev local
	Env string
}
type httpSetting struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	Throttle     string
}
type databaseSetting struct {
	Type         string
	Host         string
	Port         int
	UserName     string
	Password     string
	Database     string
	TablePrefix  string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
	LogLevel     string
}
type logSetting struct {
	Level   string
	MaxSize int
}
type aliLogSetting struct {
	Endpoint     string
	AccessKey    string
	AccessSecret string
	ProjectName  string
	LogStore     string
	Topic        string
	Source       string
}
type wxSetting struct {
	AppId  string
	Secret string
}
type redisSetting struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type duiBaSetting struct {
	AppKey    string
	AppSecret string
}
type ossSetting struct {
	CdnDomain    string
	Endpoint     string
	AccessKey    string
	AccessSecret string
	BasePath     string
	Bucket       string
	Region       string
	//用于分片上传sts授权 https://help.aliyun.com/document_detail/100624.htm?spm=a2c4g.11186623.0.0.5c452cb7TmaQGN#uicontrol-c69-98p-2bv
	StsRoleArn string
}
type amqpSetting struct {
	Url string
}
type baiDuSetting struct {
	AppKey    string
	AppSecret string
}
type baiDuImageSearchSetting struct {
	AppKey    string
	AppSecret string
}
type javaConfig struct {
	JavaLoginUrl string `binding:"required"`
	JavaWhoAmi   string `binding:"required"`
}
type zhugeConfig struct {
	AppKey    string
	AppSecret string
}
type rpcSetting struct {
	Endpoints string
	Target    string
	NonBlock  bool
	Timeout   int64
}

type activityZyh struct {
	AccessKeyId     string
	AccessKeySecret string
	Domain          string
}

type sms struct {
	Account  string
	Password string
	Url      string
}

type smsMarket struct {
	Account  string
	Password string
	Url      string
}
type promSetting struct {
	Host string
	Port int
	Path string
}

func FindOaSetting(source string) wxSetting {
	switch source {
	case "mio-srv-oa":
		return Config.MioSrvOA
	case "mio-sub-oa":
		return Config.MioSubOA
	}
	return wxSetting{}
}
