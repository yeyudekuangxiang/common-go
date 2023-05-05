
func new{{.upperStartCamelObject}}Model(db *gorm.DB,c cache.CacheConf,opts ...modelOption) *default{{.upperStartCamelObject}}Model {
	opt := initMobileOptions(Options{}, opts)
	return &default{{.upperStartCamelObject}}Model{
		cache: cache.New(c,singleFlights, stats,gorm.ErrRecordNotFound),
		cacheConf: c,
		options: opt,
		db:    db,
	}
}
func (m *default{{.upperStartCamelObject}}Model)initOptions(db *gorm.DB, opts []option) (*gorm.DB,Options){
    options:=m.options
    for _, opt := range opts {
		db, options = opt(db, options)
	}
	return db, options
}

// deleteCache 删除缓存,并且标记key已更新,10秒内从主库中读取
func (m *default{{.upperStartCamelObject}}Model) deleteCache(ctx context.Context, keys ...string) error {
	err := m.cache.DelCtx(ctx, keys...)
	if err != nil {
		return err
	}

	for _, k := range keys {
		val := ""
		err = m.cache.SetWithExpireCtx(ctx, updateTagKey+k, val, cacheUpdateSync)
		if err != nil {
			logx.Errorf("标记更新失败 %+v %+v\n", keys, err)
			cache.AddCleanTask(func() error {
				return m.cache.SetWithExpire(updateTagKey+k, val, cacheUpdateSync)
			})
		}
	}
	return nil
}
// autoDB 根据是否有更新标记选择主库还是从库
func (m *default{{.upperStartCamelObject}}Model) autoDB(ctx context.Context, key string) *gorm.DB {
	v := ""
	err := m.cache.GetCtx(ctx, updateTagKey+key, &v)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return m.db.Clauses(dbresolver.Read)
		}
		logx.Errorf("查询tag失败 %+v", err)
	}
	return m.db.Clauses(dbresolver.Write)
}
