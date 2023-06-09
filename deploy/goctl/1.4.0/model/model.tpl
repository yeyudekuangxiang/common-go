package {{.pkg}}
{{if .withCache}}
import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"context"
)
{{else}}
import (
    "github.com/zeromicro/go-zero/core/stores/cache"
    "gorm.io/gorm"
    "gorm.io/plugin/dbresolver"
    "context"
)
{{end}}
var _ {{.upperStartCamelObject}}Model = (*custom{{.upperStartCamelObject}}Model)(nil)

type (
	// {{.upperStartCamelObject}}Model is an interface to be customized, add more methods here,
	// and implement the added methods in custom{{.upperStartCamelObject}}Model.
	{{.upperStartCamelObject}}Model interface {
		{{.lowerStartCamelObject}}Model
		FindOne{{.upperStartCamelObject}}(ctx context.Context,param FindOne{{.upperStartCamelObject}}Param,opts ...option) (*{{.upperStartCamelObject}},bool,error)
		List(ctx context.Context, param List{{.upperStartCamelObject}}Param,opts ...option) ([]{{.upperStartCamelObject}},int64, error)
	    // Policy 设置从主库还是从库读 仅对customUserModel下面的方法生效
	    Policy(operation dbresolver.Operation) {{.upperStartCamelObject}}Model
	}

	custom{{.upperStartCamelObject}}Model struct {
		*default{{.upperStartCamelObject}}Model
	}
)


// New{{.upperStartCamelObject}}Model returns a model for the database table.
func New{{.upperStartCamelObject}}Model(db *gorm.DB,c cache.CacheConf, opts ...modelOption) {{.upperStartCamelObject}}Model {
	return &custom{{.upperStartCamelObject}}Model{
		default{{.upperStartCamelObject}}Model: new{{.upperStartCamelObject}}Model(db,c,opts...),
	}
}

// Policy 设置从主库还是从库读 仅对customUserModel中的方法生效
func (c *custom{{.upperStartCamelObject}}Model) Policy(operation dbresolver.Operation) {{.upperStartCamelObject}}Model {
	db := c.db.Clauses(operation).Session(&gorm.Session{})
	return New{{.upperStartCamelObject}}Model(db,c.cacheConf)
}

// FindOne{{.upperStartCamelObject}} 根据条件查询数据
func (c *custom{{.upperStartCamelObject}}Model) FindOne{{.upperStartCamelObject}}(ctx context.Context,param FindOne{{.upperStartCamelObject}}Param, opts ...option) (*{{.upperStartCamelObject}},bool,error) {
	db := c.db.WithContext(ctx)
	db ,_ = initOptions(db,c.options, opts)

	db = init{{.upperStartCamelObject}}OrderBy(db, param.OrderBy)
    //在此处组装sql


	var data {{.upperStartCamelObject}}
	err := db.Take(&data).Error

	if err == nil {
		return &data, true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return nil, false, err
}

// List 根据条件查询分页数据
func (c *custom{{.upperStartCamelObject}}Model) List(ctx context.Context, param List{{.upperStartCamelObject}}Param,opts ...option) ([]{{.upperStartCamelObject}}, int64, error) {

	db := c.db.WithContext(ctx)
	db = init{{.upperStartCamelObject}}OrderBy(db, param.OrderBy)

	//在此处组装sql


	//查询总条数
    var count int64
    if param.Limit != nil || param.Offset != nil {
        err := db.Model({{.upperStartCamelObject}}{}).Count(&count).Error
    	if err != nil {
    		return nil, 0, err
    	}
    }

    //查询列表
	db ,_ = initOptions(db,c.options, opts)
	if param.Limit != nil{
	    db = db.Limit(*param.Limit)
	}
	if param.Offset != nil{
	    db = db.Offset(*param.Offset)
	}
	list := make([]{{.upperStartCamelObject}}, 0)
	err := db.Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	if count==0{
	    count = int64(len(list))
	}
	return list, count, nil
}

// {{.upperStartCamelObject}}OrderByList {{.upperStartCamelObject}}排序列表
type {{.upperStartCamelObject}}OrderByList []{{.upperStartCamelObject}}OrderBy
// {{.upperStartCamelObject}}OrderBy 排序
type {{.upperStartCamelObject}}OrderBy struct {
	OrderBy string
}

var (
    // {{.upperStartCamelObject}}OrderBy{{.upperStartCamelPrimaryKey}}Desc 根据主键递减排序
    {{.upperStartCamelObject}}OrderBy{{.upperStartCamelPrimaryKey}}Desc = {{.upperStartCamelObject}}OrderBy{OrderBy: "{{.originPrimaryKey}} desc"}
    // {{.upperStartCamelObject}}OrderBy{{.upperStartCamelPrimaryKey}}Asc 根据主键递增排序
    {{.upperStartCamelObject}}OrderBy{{.upperStartCamelPrimaryKey}}Asc = {{.upperStartCamelObject}}OrderBy{OrderBy: "{{.originPrimaryKey}} asc"}
    {{if .hasCreatedAt}}
    // {{.upperStartCamelObject}}OrderByCreatedAtDesc 根据创建时间递减排序
    {{.upperStartCamelObject}}OrderByCreatedAtDesc = {{.upperStartCamelObject}}OrderBy{OrderBy: "created_at desc"}
    // {{.upperStartCamelObject}}OrderByCreatedAtAsc 根据创建时间递增排序
    {{.upperStartCamelObject}}OrderByCreatedAtAsc = {{.upperStartCamelObject}}OrderBy{OrderBy: "created_at asc"}
    {{end}}
    {{if .hasUpdatedAt}}
    // {{.upperStartCamelObject}}OrderByUpdatedAtDesc 根据更新时间递减排序
    {{.upperStartCamelObject}}OrderByUpdatedAtDesc = {{.upperStartCamelObject}}OrderBy{OrderBy: "updated_at desc"}
    // {{.upperStartCamelObject}}OrderByUpdatedAtAsc 根据更新时间递增排序
    {{.upperStartCamelObject}}OrderByUpdatedAtAsc = {{.upperStartCamelObject}}OrderBy{OrderBy: "updated_at asc"}
    {{end}}
)

// init{{.upperStartCamelObject}}OrderBy
func init{{.upperStartCamelObject}}OrderBy(db *gorm.DB, orderByList {{.upperStartCamelObject}}OrderByList) *gorm.DB {
	for _, item := range orderByList {
		db = db.Order(item.OrderBy)
	}
	return db
}

type FindOne{{.upperStartCamelObject}}Param struct {
    OrderBy {{.upperStartCamelObject}}OrderByList
}

type List{{.upperStartCamelObject}}Param struct {
    Limit *int
    Offset *int
    OrderBy {{.upperStartCamelObject}}OrderByList
}


