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
type ExportPointRecordListFrom struct {
	UserId    int64                       `json:"userId" form:"userId" binding:"gte=0" alias:"用户ID"`
	Nickname  string                      `json:"nickname" form:"nickname" binding:"lte=30" alias:"用户昵称"`
	OpenId    string                      `json:"openId" form:"openId" binding:"lte=50" alias:"openId"`
	Phone     string                      `json:"phone" form:"phone" binding:"lte=30" alias:"手机号"`
	StartTime time.Time                   `json:"startTime" form:"startTime" alias:"开始时间" time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime   time.Time                   `json:"endTime" form:"endTime" alias:"结束时间" time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	Type      entity.PointTransactionType `json:"type" form:"type" alias:"积分类型"`
}
type GetFileExportPageListForm struct {
	AdminId        int64                   `json:"adminId" form:"adminId" alias:"adminId"`
	Type           entity.FileExportType   `json:"type" form:"type" alias:"type"`
	Status         entity.FileExportStatus `json:"status" form:"status" alias:"status"`
	StartCreatedAt time.Time               `json:"startCreatedAt" form:"startCreatedAt" alias:"startCreatedAt" time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	EndCreatedAt   time.Time               `json:"endCreatedAt" form:"endCreatedAt" alias:"endCreatedAt" time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	controller.PageFrom
}
type GetAdjustRecordPageListForm struct {
	OpenId string                      `json:"openId" form:"openId" alias:"openId"`
	Phone  string                      `json:"phone" form:"phone" alias:"phone"`
	Type   entity.PointTransactionType `json:"type" form:"type" alias:"type"`
	controller.PageFrom
}
type AdjustUserPointForm struct {
	OpenId string                      `json:"openId" form:"openId" binding:"required" alias:"openId"`
	Phone  string                      `json:"phone" form:"phone" binding:"required" alias:"手机号"`
	Type   entity.PointTransactionType `json:"type" form:"type" binding:"oneof=SYSTEM_REDUCE SYSTEM_ADD" alias:"变动类型"`
	Value  int                         `json:"value" form:"value" binding:"gt=0" alias:"变动积分数量"`
	Note   string                      `json:"note" form:"note" binding:"gt=0,lte=200" alias:"操作备注"`
}
type AdminLoginForm struct {
	Account  string `json:"account" form:"account" binding:"required" alias:"账号"`
	Password string `json:"password" form:"password" binding:"required" alias:"密码"`
}
