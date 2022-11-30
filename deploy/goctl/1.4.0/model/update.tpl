
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data,exist, err:=m.FindOne(ctx, newData.{{.upperStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}
	if !exist{
	    return nil
	}

{{end}}	{{.keys}}
    {{if .containsIndexCache}}err{{else}}err:{{end}}= m.db.WithContext(ctx).Save(data).Error
    if err != nil {
        return err
    }
    return m.cache.DelCtx(ctx, {{.keyValues}})
    {{else}}
    return m.db.WithContext(ctx).Save({{if .containsIndexCache}}newData{{else}}data{{end}}).Error
    {{end}}
}

// UpdateColumn 更新一列数据
// {{.lowerStartCamelPrimaryKey}} 主键
// column 列名
// val 列值
// skipHook 是否跳过 Hook 方法且不追踪更新时间 true跳过 false不跳过
func (m *default{{.upperStartCamelObject}}Model) UpdateColumn(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.primaryKeyDataType}}, column string, val interface{}, skipHook bool) error {
	var err error
	{{if .withCache}}
	data, exist, err := m.FindOne(ctx, {{.lowerStartCamelPrimaryKey}})
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	{{.keys}}
	{{end}}

	if skipHook {
		err = m.db.WithContext(ctx).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).UpdateColumn(column, val).Error
	} else {
		err = m.db.WithContext(ctx).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Update(column, val).Error
	}
	if err != nil {
		return err
	}
	{{if .withCache}}
    return m.cache.DelCtx(ctx, {{.keyValues}})
    {{else}}
    return nil
    {{end}}
}

// UpdateColumns 更新多列数据
// {{.lowerStartCamelPrimaryKey}} 主键
// values map或者struct 当使用 struct 更新时，默认情况下，GORM 只会更新非零值的字段
// skipHook 是否跳过 Hook 方法且不追踪更新时间 true跳过 false不跳过
func (m *default{{.upperStartCamelObject}}Model) UpdateColumns(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.primaryKeyDataType}}, values interface{}, skipHook bool) error {
	var err error
	{{if .withCache}}
	data, exist, err := m.FindOne(ctx, {{.lowerStartCamelPrimaryKey}})
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	{{.keys}}
	{{end}}

	if skipHook {
		err = m.db.WithContext(ctx).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).UpdateColumns(values).Error
	} else {
		err = m.db.WithContext(ctx).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Updates(values).Error
	}
	if err != nil {
		return err
	}

	{{if .withCache}}
    return m.cache.DelCtx(ctx, {{.keyValues}})
    {{else}}
    return nil
    {{end}}
}