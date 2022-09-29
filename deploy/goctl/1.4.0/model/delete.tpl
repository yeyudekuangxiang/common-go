
func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data,exist, err:=m.FindOne(ctx, {{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}
	if !exist{
	    return nil
    }

{{end}}	{{.keys}}
    err {{if .containsIndexCache}}={{else}}:={{end}} m.db.WithContext(ctx).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{}).Error
    if err != nil {
        return err
    }
	err = m.cache.DelCtx(ctx, {{.keyValues}})
	return err
    {{else}}return m.db.WithContext(ctx).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{}).Error{{end}}
}