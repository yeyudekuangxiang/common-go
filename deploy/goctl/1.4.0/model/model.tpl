package {{.pkg}}
{{if .withCache}}
import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
	"context"
)
{{else}}
import (
    "gorm.io/gorm"
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
		List(ctx context.Context, param List{{.upperStartCamelObject}}Param,opts ...option) ([]{{.upperStartCamelObject}}, error)
		Page(ctx context.Context, param Page{{.upperStartCamelObject}}Param,opts ...option) ([]{{.upperStartCamelObject}}, int64, error)
	}

	custom{{.upperStartCamelObject}}Model struct {
		*default{{.upperStartCamelObject}}Model
	}
)

// New{{.upperStartCamelObject}}Model returns a model for the database table.
func New{{.upperStartCamelObject}}Model(db *gorm.DB,{{if .withCache}} c cache.CacheConf{{end}}) {{.upperStartCamelObject}}Model {
	return &custom{{.upperStartCamelObject}}Model{
		default{{.upperStartCamelObject}}Model: new{{.upperStartCamelObject}}Model(db,{{if .withCache}} c{{end}}),
	}
}

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

func (c *custom{{.upperStartCamelObject}}Model) Page(ctx context.Context, param Page{{.upperStartCamelObject}}Param,opts ...option) ([]{{.upperStartCamelObject}}, int64, error) {

	db := c.db.WithContext(ctx)
	db = init{{.upperStartCamelObject}}OrderBy(db, param.OrderBy)

	//在此处组装sql

    var count int64
	list := make([]{{.upperStartCamelObject}}, 0)
	err := db.Model({{.upperStartCamelObject}}{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	db ,_ = initOptions(db,c.options, opts)
	err = db.Limit(param.Limit).Offset(param.Offset).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
func (c *custom{{.upperStartCamelObject}}Model) List(ctx context.Context, param List{{.upperStartCamelObject}}Param,opts ...option) ([]{{.upperStartCamelObject}}, error) {

	db := c.db.WithContext(ctx)
	db ,_ = initOptions(db,c.options, opts)
	db = init{{.upperStartCamelObject}}OrderBy(db, param.OrderBy)

	//在此处组装sql


	list := make([]{{.upperStartCamelObject}}, 0)
	err := db.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

type {{.upperStartCamelObject}}OrderByList []{{.upperStartCamelObject}}OrderBy
type {{.upperStartCamelObject}}OrderBy struct {
	OrderBy string
}

var (
    {{.upperStartCamelObject}}OrderBy{{.upperStartCamelPrimaryKey}}Desc = {{.upperStartCamelObject}}OrderBy{OrderBy: "{{.originPrimaryKey}} desc"}
    {{.upperStartCamelObject}}OrderBy{{.upperStartCamelPrimaryKey}}Asc = {{.upperStartCamelObject}}OrderBy{OrderBy: "{{.originPrimaryKey}} asc"}
    {{if .hasCreatedAt}}
     {{.upperStartCamelObject}}OrderByCreatedAtDesc = {{.upperStartCamelObject}}OrderBy{OrderBy: "created_at desc"}
        {{.upperStartCamelObject}}OrderByCreatedAtAsc = {{.upperStartCamelObject}}OrderBy{OrderBy: "created_at asc"}
    {{end}}
    {{if .hasUpdatedAt}}
         {{.upperStartCamelObject}}OrderByUpdatedAtDesc = {{.upperStartCamelObject}}OrderBy{OrderBy: "updated_at desc"}
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

type Page{{.upperStartCamelObject}}Param struct {
    Limit int
    Offset int
    OrderBy {{.upperStartCamelObject}}OrderByList
}

type List{{.upperStartCamelObject}}Param struct {
    OrderBy {{.upperStartCamelObject}}OrderByList
}

