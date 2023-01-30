package admin

import (
	"mio/internal/app/mp2c/controller"
	"mio/internal/pkg/model/entity"
	"time"
)

type GetUserForm struct {
	Id int64
}

type IDForm struct {
	ID int64 `json:"id" form:"id" alias:"id" binding:"required,gte=1" `
}

type GetPointRecordPageListFrom struct {
	UserId          int64                       `json:"userId" form:"userId" binding:"gte=0" alias:"用户ID"`
	Nickname        string                      `json:"nickname" form:"nickname" binding:"lte=30" alias:"用户昵称"`
	OpenId          string                      `json:"openId" form:"openId" binding:"lte=50" alias:"openId"`
	Phone           string                      `json:"phone" form:"phone" binding:"lte=30" alias:"手机号"`
	StartTime       time.Time                   `json:"startTime" form:"startTime" alias:"开始时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime         time.Time                   `json:"endTime" form:"endTime" alias:"结束时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
	Type            entity.PointTransactionType `json:"type" form:"type" alias:"积分类型"`
	StartExpireTime string                      `json:"startExpireTime" form:"startExpireTime" binding:"omitempty,datetime=2006-01-02"`
	EndExpireTime   string                      `json:"endExpireTime" form:"endExpireTime" binding:"omitempty,datetime=2006-01-02"`
	controller.PageFrom
}
type ExportPointRecordListFrom struct {
	UserId          int64                       `json:"userId" form:"userId" binding:"gte=0" alias:"用户ID"`
	Nickname        string                      `json:"nickname" form:"nickname" binding:"lte=30" alias:"用户昵称"`
	OpenId          string                      `json:"openId" form:"openId" binding:"lte=50" alias:"openId"`
	Phone           string                      `json:"phone" form:"phone" binding:"lte=30" alias:"手机号"`
	StartTime       time.Time                   `json:"startTime" form:"startTime" alias:"开始时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime         time.Time                   `json:"endTime" form:"endTime" alias:"结束时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
	StartExpireTime string                      `json:"startExpireTime" form:"startExpireTime" binding:"omitempty,datetime=2006-01-02"`
	EndExpireTime   string                      `json:"endExpireTime" form:"endExpireTime" binding:"omitempty,datetime=2006-01-02"`
	Type            entity.PointTransactionType `json:"type" form:"type" alias:"积分类型"`
}
type GetFileExportPageListForm struct {
	AdminId        int64                   `json:"adminId" form:"adminId" alias:"adminId"`
	Type           int                     `json:"type" form:"type" alias:"type"`
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
	ID        int64     `json:"id" form:"id" alias:"用户id" binding:"gte=0"`
	Mobile    string    `json:"mobile" form:"mobile" alias:"手机号码"`
	Nickname  string    `json:"nickname" form:"nickname" alias:"昵称"`
	State     int       `json:"state" form:"state" alias:"状态" binding:"min=0,max=3"`
	Position  string    `json:"position" form:"position" alias:"身份"`
	Partners  int       `json:"partners" form:"partners" alias:"合作伙伴"`
	Auth      int       `json:"auth"`
	StartTime time.Time `json:"startTime" form:"startTime"  time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime   time.Time `json:"endTime" form:"endTime"  time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	controller.PageFrom
}

/*common start*/
type CommentDeleteRequest struct {
	Reason string `json:"reason" form:"reason" alias:"reason"`
	IDForm
}

/*common end*/

/*topic start*/
type TopicListRequest struct {
	ID         int64  `json:"id" form:"id" alias:"帖子id" binding:"gte=0"`
	Title      string `json:"title" form:"title" alias:"帖子标题"`
	TagId      string `json:"tagId" form:"tagId" alias:"标签id" binding:"gte=0"`
	UserId     int64  `json:"userId" form:"userId" alias:"用户id" binding:"gte=0"`
	UserName   string `json:"userName" form:"userName" alias:"用户名"`
	Status     int    `json:"status" form:"status" alias:"审核状态" binding:"min=0,max=4"`
	IsTop      int    `json:"isTop" form:"isTop" alias:"是否置顶" binding:"min=0,max=1"`
	IsEssence  int    `json:"isEssence" form:"isEssence" alias:"是否精华" binding:"min=0,max=1"`
	IsPartners int    `json:"isPartners" form:"isPartners" alias:"乐活家" binding:"min=0"`
	Position   string `json:"position" form:"position" alias:"身份"`
	controller.PageFrom
}

type CreateTopicRequest struct {
	Title   string   `json:"title" form:"title" alias:"title" binding:"required,min=2,max=64"`
	Content string   `json:"content" form:"content" alias:"content" binding:"min=0,max=10000"`
	Images  []string `json:"images" form:"images" alias:"images" binding:"required,min=1,max=12"`
	TagIds  []int64  `json:"tagIds" form:"tagIds" alias:"tagIds" binding:"min=0,max=2"`
}

type UpdateTopicRequest struct {
	ID int64 `json:"id" form:"id" alias:"id" binding:"required,gte=1"`
	CreateTopicRequest
}

type ChangeTopicStatus struct {
	ID        int64  `json:"id" form:"id" alias:"id" binding:"required,gte=1"`
	IsTop     int    `json:"isTop" form:"isTop" alias:"isTop" binding:"oneof=0 1"`
	IsEssence int    `json:"isEssence" form:"isEssence" alias:"isEssence" binding:"oneof=0 1"`
	Status    int    `json:"status" form:"status" alias:"审核状态" binding:"min=0,max=4"`
	Reason    string `json:"reason" form:"reason" alias:"审核未通过理由"`
}

/*topic end*/

/*comment start*/
type CommentListRequest struct {
	Comment string `json:"comment" form:"comment" alias:"comment"`
	UserId  int64  `json:"userId" form:"userId" alias:"userId" binding:"gte=0"`
	TopicId int64  `json:"topicId" form:"topicId" alias:"topicId" binding:"gte=0"`
	controller.PageFrom
}

/*comment end*/

/*tag start*/
type TagListRequest struct {
	Name        string `json:"name" form:"name" alias:"name"`
	Description string `json:"description" form:"description" alias:"description"`
	controller.PageFrom
}
type CreateTagRequest struct {
	Name        string `json:"name" form:"name" alias:"name" binding:"required,min=2,max=64"`
	Description string `json:"description" form:"description" alias:"description" binding:"min=0,max=200"`
	Image       string `json:"image" form:"image" alias:"image"`
}
type UpdateTagRequest struct {
	ID int64 `json:"id" form:"id" alias:"id" binding:"required,gte=1"`
	CreateTagRequest
}

/*tag end*/

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
	Type     entity.BannerType   `json:"type" form:"type" binding:"oneof=mini path h5" alias:"跳转类型"`
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
	Type     entity.BannerType   `json:"type" form:"type" binding:"oneof=mini path h5" alias:"跳转类型"`
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

type CreateDuiBaActivityForm struct {
	Name        string                      `json:"name" form:"name" binding:"required" alias:"链接名称"`
	Cid         int64                       `json:"cid" form:"cid" binding:"" alias:"渠道id"`
	Type        entity.DuiBaActivityType    `json:"type" form:"type" binding:"oneof=1 2" alias:"类别"`
	IsShare     entity.DuiBaActivityIsShare `json:"isShare" form:"isShare" binding:"oneof=1 2" alias:"是否分享"`
	IsPhone     entity.DuiBaActivityIsPhone `json:"isPhone" form:"isPhone" binding:"oneof=1 2" alias:"是否绑定手机号"`
	ActivityUrl string                      `json:"activityUrl" form:"activityUrl" binding:"required" alias:"链接"`
	ActivityId  string                      `json:"activityId" form:"activityId" binding:"required" alias:"活动"`
	RiskLimit   int                         `json:"riskLimit" form:"riskLimit" binding:"oneof=0 1 2 3 4" alias:"允许的风险等级"`
}

type UpdateDuiBaActivityForm struct {
	Id          int64                       `json:"id" form:"id" binding:"required" alias:"id"`
	Name        string                      `json:"name" form:"name" binding:"required" alias:"链接名称"`
	Cid         int64                       `json:"cid" form:"cid" binding:"" alias:"渠道id"`
	Type        entity.DuiBaActivityType    `json:"type" form:"type" binding:"oneof=1 2" alias:"类别"`
	IsShare     entity.DuiBaActivityIsShare `json:"isShare" form:"isShare" binding:"oneof=1 2" alias:"是否分享"`
	IsPhone     entity.DuiBaActivityIsPhone `json:"isPhone" form:"isPhone" binding:"oneof=1 2" alias:"是否绑定手机号"`
	ActivityUrl string                      `json:"activityUrl" form:"activityUrl" binding:"required" alias:"链接"`
	ActivityId  string                      `json:"activityId" form:"activityId" binding:"required" alias:"活动"`
	RiskLimit   int                         `json:"riskLimit" form:"riskLimit" binding:"oneof=0 1 2 3 4" alias:"允许的风险等级"`
}

type GetDuiBaActivityPageForm struct {
	Cid        int64                    `json:"cid" form:"cid" binding:"" alias:"渠道id"`
	Type       entity.DuiBaActivityType `json:"type" form:"type" binding:"" alias:"类别"`
	ActivityId string                   `json:"activityId" form:"activityId" binding:"" alias:"活动"`
	Name       string                   `json:"name" form:"name" binding:"" alias:"链接名称"`
	controller.PageFrom
}

type DeleteDuiBaActivityForm struct {
	Id int64 `json:"id" form:"id" binding:"required" alias:"id"`
}

type ShowDuiBaActivityForm struct {
	Id int64 `json:"id" form:"id" binding:"required" alias:"id"`
}

//user start

type UpdateUser struct {
	IDForm
	Status   int    `json:"status,omitempty" form:"status" alias:"status"`
	Position string `json:"position,omitempty" form:"position" alias:"position"`
	Partners int    `json:"partners,omitempty" form:"partners" alias:"partners"`
	Auth     int    `json:"auth,omitempty" form:"auth" alias:"auth"`
}

//user end

type SendCouponReq struct {
	OpenId string  `json:"openId" form:"openId" binding:"required"`
	Ch     string  `json:"ch" form:"ch"`
	Amount float64 `json:"amount" form:"amount"`
}

type UpdateUserRisk struct {
	Ids  string `json:"ids" form:"ids"    binding:"required"`
	Type int    `json:"type" form:"type"  binding:"oneof=1 2 3" alias:"提交类型  1 uid、2 手机号、3 openid "`
	Risk int    `json:"risk" form:"risk"  binding:"oneof=0 1 2 3 4" alias:"允许的风险等级"`
}

//index_icon start 金刚位

type CreateIndexIconForm struct {
	Title  string                 `json:"title" form:"title" binding:"required" alias:"标题"`
	RowNum string                 `json:"rowNum" form:"rowNum" binding:"oneof=1 2" alias:"行"`
	Sort   int8                   `json:"sort" form:"sort" binding:"" alias:"排序"`
	Status entity.IndexIconStatus `json:"status" form:"status" binding:"oneof=1 2" alias:"是否删除 1 正常 2删除"`
	IsOpen int8                   `json:"isOpen" form:"isOpen" binding:"oneof=1 2" alias:"是否开启 1开启 2关闭"`
	Pic    string                 `json:"pic" form:"pic" binding:"" alias:"图标"`
}

type UpdateIndexIconForm struct {
	Id     int64                  `json:"id" form:"id" binding:"required" alias:"id"`
	Title  string                 `json:"title" form:"title" binding:"required" alias:"标题"`
	Type   string                 `json:"type" form:"type" binding:"" alias:"类型"`
	RowNum string                 `json:"rowNum" form:"rowNum" binding:"oneof=1 2" alias:"行"`
	Sort   int8                   `json:"sort" form:"sort" binding:"" alias:"排序"`
	Status entity.IndexIconStatus `json:"status" form:"status" binding:"oneof=1 2" alias:"是否删除 1 正常 2删除"`
	IsOpen int8                   `json:"isOpen" form:"isOpen" binding:"oneof=1 2" alias:"是否开启 1开启 2关闭"`
	Pic    string                 `json:"pic" form:"pic" binding:"" alias:"图标"`
}

type GetIndexIconPageForm struct {
	Title  string                 `json:"title" form:"title" binding:"" alias:"标题"`
	IsOpen entity.IndexIconIsOpen `json:"isOpen" form:"isOpen" binding:"" alias:"是否开启 1开启 2关闭"`
	controller.PageFrom
}

type DeleteIndexIconForm struct {
	Id int64 `json:"id" form:"id" binding:"required" alias:"id"`
}

type ShowIndexIconForm struct {
	Id int64 `json:"id" form:"id" binding:"required" alias:"id"`
}

type PointTransactionTypeInfo struct {
	Type     string `json:"type"`
	TypeText string `json:"typeText"`
}
