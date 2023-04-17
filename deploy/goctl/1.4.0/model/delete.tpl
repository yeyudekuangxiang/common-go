// Delete 根据主键删除数据
func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .containsIndexCache}}
	//查询数据用于删除索引
	data,exist, err:=m.FindOne(ctx, {{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}
	if !exist{
	    return nil
    }

{{end}}	{{.keys}}
    //删除数据
    err {{if .containsIndexCache}}={{else}}:={{end}} m.db.WithContext(ctx).Clauses(dbresolver.Write).Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{}).Error
    if err != nil {
        return err
    }
    //删除缓存,标记删除
	err = m.deleteCache(ctx, {{.keyValues}})
	return err
}