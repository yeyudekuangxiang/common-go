// Insert 创建数据
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}) error {

	err := m.db.WithContext(ctx).Clauses(dbresolver.Write).Create(data).Error
	if err != nil {
		return err
	}

	{{.keys}}
    //删除缓存,标记删除
	return m.deleteCache(ctx, {{.keyValues}})
}
// InsertBatch 批量创建数据
func (m *default{{.upperStartCamelObject}}Model) InsertBatch(ctx context.Context, data *[]{{.upperStartCamelObject}}, batchSize int) error {
	err := m.db.WithContext(ctx).Clauses(dbresolver.Write).CreateInBatches(data, batchSize).Error
	if err!=nil{
	    return err
	}
	keys:=make([]string,0)
	for _,d:=range *data{
	    keys = append(keys,m.formatKeys(d)...)
	}
	//删除缓存,标记删除
    return m.deleteCache(ctx,keys...)
}
// formatKeys 格式化主键和索引的缓存key
func (m *default{{.upperStartCamelObject}}Model) formatKeys(data {{.upperStartCamelObject}})[]string {
    {{.keys}}
    return []string{ {{.keyValues}} }
}