
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data,exist, err:=m.FindOne(ctx, newData.{{.upperStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}

{{end}}	{{.keys}}
    {{if .containsIndexCache}}err{{else}}err:{{end}}= m.db.WithContext(ctx).Save(data).Error
    if err != nil {
        return err
    }
    return m.cache.DelCtx(ctx, publicCouponCardCouponCardIdKey)
    {{else}}
    return m.db.WithContext(ctx).Save(data).Error
    {{end}}
}
