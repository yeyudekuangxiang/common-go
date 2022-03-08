package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iris-contrib/schema"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/date"
	"strconv"
	"strings"
	"time"
)
import (
	"errors"
)

var (
	decoder = schema.NewDecoder() // form, url, schema.
)

func init() {
	decoder.AddAliasTag("form", "json")
	decoder.ZeroEmpty(true)
}

// param error
func paramError(name string) error {
	return errors.New(fmt.Sprintf("unable to find param value '%s'", name))
}

func PostForm(ctx *gin.Context, name string) string {
	return ctx.PostForm(name)
}

func PostFormRequired(ctx *gin.Context, name string) (string, error) {
	str := PostForm(ctx, name)
	if len(str) == 0 {
		return "", errors.New("参数：" + name + "不能为空")
	}
	return str, nil
}

func PostFormDefault(ctx *gin.Context, name, def string) string {
	return ctx.DefaultPostForm(name, def)
}

func PostFormInt(ctx *gin.Context, name string) (int, error) {
	str := ctx.PostForm(name)
	if str == "" {
		return 0, paramError(name)
	}
	return strconv.Atoi(str)
}

func PostFormIntDefault(ctx *gin.Context, name string, def int) int {
	if v, err := PostFormInt(ctx, name); err == nil {
		return v
	}
	return def
}

func PostFormInt64(ctx *gin.Context, name string) (int64, error) {
	str := ctx.PostForm(name)
	if str == "" {
		return 0, paramError(name)
	}
	return strconv.ParseInt(str, 10, 64)
}

func PostFormInt64Default(ctx *gin.Context, name string, def int64) int64 {
	if v, err := PostFormInt64(ctx, name); err == nil {
		return v
	}
	return def
}

func PostFormInt64Array(ctx *gin.Context, name string) []int64 {
	str := ctx.PostForm(name)
	if str == "" {
		return nil
	}
	ss := strings.Split(str, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []int64
	for _, v := range ss {
		item, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		ret = append(ret, item)
	}
	return ret
}

func PostFormStringArray(ctx *gin.Context, name string) []string {
	str := ctx.PostForm(name)
	if len(str) == 0 {
		return nil
	}
	ss := strings.Split(str, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		ret = append(ret, s)
	}
	return ret
}

func PostFormBool(ctx *gin.Context, name string) (bool, error) {
	str := ctx.PostForm(name)
	if str == "" {
		return false, paramError(name)
	}
	return strconv.ParseBool(str)
}

// 从请求中获取日期
func FormDate(ctx *gin.Context, name string) *time.Time {
	value := PostForm(ctx, name)
	if simple.IsBlank(value) {
		return nil
	}
	layouts := []string{date.FmtDateTime, date.FmtDate, date.FmtDateTimeNoSeconds}
	for _, layout := range layouts {
		if ret, err := date.Parse(value, layout); err == nil {
			return &ret
		}
	}
	return nil
}

func GetPaging(ctx *gin.Context) *simple.Paging {
	page := PostFormIntDefault(ctx, "page", 1)
	limit := PostFormIntDefault(ctx, "limit", 20)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &simple.Paging{Page: page, Limit: limit}
}
