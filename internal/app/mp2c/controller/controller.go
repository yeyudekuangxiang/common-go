package controller

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
