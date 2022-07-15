package controller

import "github.com/shopspring/decimal"

type PageFrom struct {
	Page     int `json:"page" form:"page" binding:"gt=0" alias:"页码"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"gt=0" alias:"每页数量"`
}

func (p PageFrom) Limit() int {
	return p.PageSize
}
func (p PageFrom) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func NewPageResult(list interface{}, total int64, page PageFrom) PageResult {
	totalPage := int(decimal.NewFromInt(total).Div(decimal.NewFromInt(int64(page.PageSize))).Ceil().IntPart())
	return PageResult{
		List:      list,
		Page:      page.Page,
		PageSize:  page.PageSize,
		Total:     total,
		TotalPage: totalPage,
		HasNext:   page.Page < totalPage,
	}
}

type PageResult struct {
	List      interface{} `json:"list"`
	Page      int         `json:"page"`
	PageSize  int         `json:"pageSize"`
	Total     int64       `json:"total"`
	TotalPage int         `json:"totalPage"`
	HasNext   bool        `json:"hasNext"`
}
