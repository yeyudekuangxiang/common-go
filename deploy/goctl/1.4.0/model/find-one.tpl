
func (m *default{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}},bool, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.cache.TakeCtx(ctx, &resp, {{.cacheKeyVariable}}, func(val interface{}) error {
		return m.db.WithContext(ctx).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Take(&resp).Error
	})

	if err == nil {
		return &resp, true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return nil, false, err
    {{else}}var resp {{.upperStartCamelObject}}
    err:=m.db.WithContext(ctx).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Take(&resp).Error
 	if err == nil {
 		return &resp, true, nil
 	}
 	if err == gorm.ErrRecordNotFound {
 		return nil, false, nil
 	}
 	return nil, false, err{{end}}
}
