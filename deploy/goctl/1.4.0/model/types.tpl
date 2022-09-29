
type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}cache cache.Cache{{end}}
		db    *gorm.DB
	}

    // Order 请手动设置主键字段  gorm:"primaryKey" 否则 default{{.upperStartCamelObject}}Model.Update 方法会抛异常
	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
