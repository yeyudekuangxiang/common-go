func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}(ctx context.Context, {{.in}}, opts ...option) (*{{.upperStartCamelObject}},bool, error) {
	{{if .withCache}}
	{{.cacheKey}}
    var resp {{.upperStartCamelObject}}
	db,op:= initOptions(m.db.WithContext(ctx),m.options,opts)
    var primaryKey {{.primaryKeyDataType}}
    var found bool
    var err error

	if op.skipFindCache {
	    goto skipFindCache
	}

	//从缓存中根据索引查询主键内容
	err = m.cache.TakeWithExpireCtx(ctx, &primaryKey, {{.cacheKeyVariable}},
		func(val interface{}, expire time.Duration) (err error) {
			//缓存中没有查到主键的缓存 从数据库中根据索引查询数据
			if err := db.Model(&resp).Where("{{.originalField}}", {{.lowerStartCamelField}}).First(&resp).Error; err != nil {
				// 未查到或者有异常返回err
				return err
			}
			//查询数据信息
			primaryKey = resp.{{.upperStartCamelPrimaryKey}}
			found = true
			//将数据存到缓存里面
			return m.cache.SetWithExpireCtx(ctx, m.formatPrimary(primaryKey), resp,
				expire+cacheSafeGapBetweenIndexAndPrimary)
		})

	//判断err
	if err != nil {
		//未从数据库中查到数据
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		//查询异常
		return nil, false, err
	}

	//从数据库中查到了数据
	if found {
		return &resp, true, nil
	}

    skipFindCache:
	//从缓存中查到了 主键的索引 从缓存中查询数据
	return m.FindOne(ctx, primaryKey,opts...)
	{{else}}
	var resp {{.upperStartCamelObject}}
    db,_:= initOptions(m.db.WithContext(ctx),m.options,opts)
    err := db.Model(&resp).Where("{{.originalField}}", {{.lowerStartCamelField}}).First(&resp).Error
    if err==nil{
        return &resp,true,nil
    }
    if err == gorm.ErrRecordNotFound {
        return nil,false,nil
    }
    return nil,false,err
	{{end}}
}
