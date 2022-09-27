
func new{{.upperStartCamelObject}}Model(db *gorm.DB,{{if .withCache}} c cache.CacheConf{{end}}) *default{{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		{{if .withCache}}cache: cache.New(c,nil,&cache.Stat{},gorm.ErrRecordNotFound),{{end}}
		db:    db,
	}
}
