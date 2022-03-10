package config

var App = app{
	App:      appSetting{},
	Http:     httpSetting{},
	Database: databaseSetting{},
	Log:      logSetting{},
	Weapp:    wxSetting{},
	Wxoa:     wxSetting{},
}

type app struct {
	App      appSetting      `ini:"app"`
	Http     httpSetting     `ini:"http"`
	Database databaseSetting `ini:"database"`
	Log      logSetting      `ini:"log"`
	Weapp    wxSetting       `ini:"weapp"`
	Wxoa     wxSetting       `ini:"wxoa"`
}
type appSetting struct {
	TokenKey string
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
type wxSetting struct {
	AppId  string
	Secret string
}
