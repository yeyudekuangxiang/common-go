// Update 更新数据
// Update会保存所有的字段，即使字段是零值
Update(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) error
// UpdateColumn 更新一列数据
// {{.lowerStartCamelPrimaryKey}} 主键
// column 列名
// val 列值
// skipHook 是否跳过 Hook 方法且不追踪更新时间 true跳过 false不跳过
UpdateColumn(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.primaryKeyDataType}}, column string, val interface{}, skipHook bool) error
// UpdateColumns 更新多列数据
// {{.lowerStartCamelPrimaryKey}} 主键
// values map或者struct 当使用 struct 更新时，默认情况下，GORM 只会更新非零值的字段
// skipHook 是否跳过 Hook 方法且不追踪更新时间 true跳过 false不跳过
UpdateColumns(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.primaryKeyDataType}}, values interface{}, skipHook bool) error
