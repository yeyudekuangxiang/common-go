
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.keys}}
	err := m.db.WithContext(ctx).Create(data).Error
	if err != nil {
		return err
	}

	return m.cache.DelCtx(ctx, {{.keyValues}})
    {{else}}
	return m.db.WithContext(ctx).Create(data).Error
    {{end}}
}
