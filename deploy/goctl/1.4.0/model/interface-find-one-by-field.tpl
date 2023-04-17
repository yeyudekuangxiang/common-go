// FindOneBy{{.upperField}} 根据索引查询数据
FindOneBy{{.upperField}}(ctx context.Context, {{.in}},opts ...option) (*{{.upperStartCamelObject}},bool, error)