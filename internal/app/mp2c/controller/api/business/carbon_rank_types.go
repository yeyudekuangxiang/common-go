package business

import "mio/internal/app/mp2c/controller"

type GetUserRankListForm struct {
	DateType string `json:"dateType" form:"dateType" binding:"oneof=day week month" alias:"榜单类型"`
	controller.PageFrom
}
