
type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}cache cache.Cache{{end}}
		db    *gorm.DB
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
