
// FindOne 根据主键查询信息
func (m *default{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}},opts ...option) (*{{.upperStartCamelObject}},bool, error) {

    {{.cacheKey}}
    db := m.autoDB(ctx,{{.cacheKeyVariable}})

	{{if .withCache}}

	db,op:= initOptions(db,&m.options,opts)

	var resp {{.upperStartCamelObject}}
	var err error

	if op.skipFindCache {
	    err = db.Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Take(&resp).Error
	    goto skipFindCache
	}

	err= m.cache.TakeCtx(ctx, &resp, {{.cacheKeyVariable}}, func(val interface{}) error {
		return db.Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Take(&resp).Error
	})

    skipFindCache:
	if err == nil {
		return &resp, true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return nil, false, err
    {{else}}
    db,_ = initOptions(db,&m.options,opts)
    var resp {{.upperStartCamelObject}}
    err:=db.Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Take(&resp).Error
 	if err == nil {
 		return &resp, true, nil
 	}
 	if err == gorm.ErrRecordNotFound {
 		return nil, false, nil
 	}
 	return nil, false, err{{end}}
}
