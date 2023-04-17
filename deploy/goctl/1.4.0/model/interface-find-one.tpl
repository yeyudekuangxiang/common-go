// FindOne 根据主键查询数据
FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}, opts ...option) (*{{.upperStartCamelObject}}, bool, error)