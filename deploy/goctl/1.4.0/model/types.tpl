
type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
		cache cache.Cache
		cacheConf cache.CacheConf
		options Options
		db    *gorm.DB
	}

    // {{.upperStartCamelObject}}
	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
