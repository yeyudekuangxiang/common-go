package initialize

func Initialize(configPath string) {
	InitIni(configPath)
	InitLog()
	InitDB()
	InitRedis()
	InitValidator()
	InitWeapp()
	initOss()
}
