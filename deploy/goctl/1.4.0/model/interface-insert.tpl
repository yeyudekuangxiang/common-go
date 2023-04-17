// Insert 创建数据
Insert(ctx context.Context, data *{{.upperStartCamelObject}}) error
// InsertBatch 批量创建数据
InsertBatch(ctx context.Context, data *[]{{.upperStartCamelObject}}, batchSize int) error