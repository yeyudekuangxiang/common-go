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
	StartTime time.Time                   `json:"startTime" form:"startTime" alias:"开始时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime   time.Time                   `json:"endTime" form:"endTime" alias:"结束时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
	Type      entity.PointTransactionType `json:"type" form:"type" alias:"积分类型"`
	controller.PageFrom
}
type ExportPointRecordListFrom struct {
	UserId    int64                       `json:"userId" form:"userId" binding:"gte=0" alias:"用户ID"`
	Nickname  string                      `json:"nickname" form:"nickname" binding:"lte=30" alias:"用户昵称"`
	OpenId    string                      `json:"openId" form:"openId" binding:"lte=50" alias:"openId"`
	Phone     string                      `json:"phone" form:"phone" binding:"lte=30" alias:"手机号"`
	StartTime time.Time                   `json:"startTime" form:"startTime" alias:"开始时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime   time.Time                   `json:"endTime" form:"endTime" alias:"结束时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
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
	UserId    int64                       `json:"userId" form:"userId" alias:"userId"`
	Nickname  string                      `json:"nickname" form:"nickname" alias:"nickname"`
	OpenId    string                      `json:"openId" form:"openId" alias:"openId"`
	Phone     string                      `json:"phone" form:"phone" alias:"phone"`
	Type      entity.PointTransactionType `json:"type" form:"type" alias:"type"`
	AdminId   int                         `json:"adminId" form:"adminId" alias:"adminId"`
	StartTime time.Time                   `json:"startTime" form:"startTime" alias:"开始时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime   time.Time                   `json:"endTime" form:"endTime" alias:"结束时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
	controller.PageFrom
}
type AdjustUserPointForm struct {
	OpenId string                      `json:"openId" form:"openId" binding:"required" alias:"openId"`
	Phone  string                      `json:"phone" form:"phone" binding:"required" alias:"手机号"`
	Type   entity.PointTransactionType `json:"type" form:"type" binding:"oneof=SYSTEM_REDUCE SYSTEM_ADD" alias:"变动类型"`
	Value  int64                       `json:"value" form:"value" binding:"gt=0" alias:"变动积分数量"`
	Note   string                      `json:"note" form:"note" binding:"gt=0,lte=200" alias:"操作备注"`
}
type AdminLoginForm struct {
	Account  string `json:"account" form:"account" binding:"required" alias:"账号"`
	Password string `json:"password" form:"password" binding:"required" alias:"密码"`
}

type UserPageListForm struct {
	Mobile    string    `json:"mobile" form:"mobile" alias:"手机号码"`
	Nickname  string    `json:"nickname" form:"nickname"`
	ID        int64     `json:"id" form:"id" binding:"gte=0" alias:"topic id"`
	StartTime time.Time `json:"startTime" form:"startTime"  time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime   time.Time `json:"endTime" form:"endTime"  time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
}

type GetUserChannelPageForm struct {
	Cid  int    `json:"cid" form:"cid" binding:"" alias:"渠道id"`
	Pid  int    `json:"pid" form:"pid" binding:"required" alias:"分类id不能为空"`
	Name string `json:"name" form:"name" binding:"" alias:"渠道名称"`
	Code string `json:"code" form:"code" binding:"" alias:"渠道code"`
	controller.PageFrom
}

type CreateUserChannelForm struct {
	Cid     int64  `json:"cid" form:"cid" binding:"gte=0" alias:"渠道id"`
	Pid     int64  `json:"pid" form:"pid" binding:"gte=0" alias:"渠道父级ID"`
	Name    string `json:"name" form:"name" binding:"required" alias:"渠道名称"`
	Code    string `json:"code" form:"code" binding:"required" alias:"渠道code"`
	Company string `json:"company" form:"company" binding:"" alias:"公司名称"`
}

type CreateBannerForm struct {
	Name     string              `json:"name" form:"name" binding:"required" alias:"banner名称"`
	ImageUrl string              `json:"imageUrl" form:"imageUrl" binding:"required" alias:"轮播图图片"`
	Scene    entity.BannerScene  `json:"scene" form:"scene" binding:"" alias:"轮播图场景"`
	Type     entity.BannerType   `json:"type" form:"type" binding:"oneof=mini path" alias:"跳转类型"`
	AppId    string              `json:"appId" form:"appId" binding:"" alias:"小程序appid(跳转到三方小程序需要)"`
	Sort     int                 `json:"sort" form:"sort" binding:"" alias:"排序"`
	Redirect string              `json:"redirect" form:"redirect" binding:"" alias:"跳转路径"`
	Status   entity.BannerStatus `json:"status" form:"status" binding:"oneof=1 2" alias:"上线和下线状态"`
}

type UpdateBannerForm struct {
	Id       int64               `json:"id" form:"id" binding:"required" alias:"bannerId"`
	Name     string              `json:"name" form:"name" binding:"required" alias:"banner名称"`
	ImageUrl string              `json:"imageUrl" form:"imageUrl" binding:"required" alias:"轮播图图片"`
	Scene    entity.BannerScene  `json:"scene" form:"scene" binding:"" alias:"轮播图场景"`
	Type     entity.BannerType   `json:"type" form:"type" binding:"oneof=mini path" alias:"跳转类型"`
	AppId    string              `json:"appId" form:"appId" binding:"" alias:"小程序appid(跳转到三方小程序需要)"`
	Sort     int                 `json:"sort" form:"sort" binding:"" alias:"排序"`
	Redirect string              `json:"redirect" form:"redirect" binding:"" alias:"跳转路径"`
	Status   entity.BannerStatus `json:"status" form:"status" binding:"oneof=1 2" alias:"上线和下线状态"`
}

type GetBannerPageForm struct {
	Name   string `json:"name" form:"name" binding:"" alias:"banner名称"`
	Scene  string `json:"scene" form:"scene" binding:"" alias:"轮播图场景"`
	Status int8   `json:"status" form:"status" binding:"" alias:"上线和下线状态"`
	controller.PageFrom
}
