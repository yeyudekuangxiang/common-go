
type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}cache cache.Cache{{end}}
		{{if .withCache}}cacheConf cache.CacheConf{{end}}
		options Options
		db    *gorm.DB
	}

    // {{.upperStartCamelObject}}
	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
