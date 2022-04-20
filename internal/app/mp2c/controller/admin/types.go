package admin

import (
	"mio/internal/app/mp2c/controller"
	"mio/internal/pkg/model/entity"
	"time"
)

type GetUserForm struct {
	Id int64
}
type GetPointRecordPageListFrom struct {
	UserId    int64                       `json:"userId" form:"userId" binding:"gte=0" alias:"用户ID"`
	Nickname  string                      `json:"nickname" form:"nickname" binding:"lte=30" alias:"用户昵称"`
	OpenId    string                      `json:"openId" form:"openId" binding:"lte=50" alias:"openId"`
	Phone     string                      `json:"phone" form:"phone" binding:"lte=30" alias:"手机号"`
	StartTime time.Time                   `json:"startTime" form:"startTime" alias:"开始时间" time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime   time.Time                   `json:"endTime" form:"endTime" alias:"结束时间" time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	Type      entity.PointTransactionType `json:"type" form:"type" alias:"积分类型"`
	controller.PageFrom
}
type GetFileExportPageListForm struct {
	Type           entity.FileExportType   `json:"type" form:"type" alias:"type"`
	Status         entity.FileExportStatus `json:"status" form:"status" alias:"status"`
	StartCreatedAt time.Time               `json:"startCreatedAt" form:"startCreatedAt" alias:"startCreatedAt" time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	EndCreatedAt   time.Time               `json:"endCreatedAt" form:"endCreatedAt" alias:"endCreatedAt" time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	controller.PageFrom
}
