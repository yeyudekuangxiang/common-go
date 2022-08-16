package entity

import (
	"mio/internal/pkg/model"
)

type Banner struct {
	ID         int64        `json:"id" gorm:"primaryKey;type:serial8;not null;comment:轮播图"`
	Name       string       `json:"name" gorm:"type:varchar(1000);not null;comment:轮播图名称"`
	ImageUrl   string       `json:"imageUrl" gorm:"type:varchar(1000);not null;comment:轮播图图片"`
	Scene      BannerScene  `json:"scene" gorm:"type:varchar(20);not null;default:'home';comment:证书场景 home首页 event携手 topic社区"`
	Type       BannerType   `json:"type" gorm:"type:varchar(20);not null;default:'path';comment:跳转类型 mini第三方小程序 path内部小程序路径 h5 h5链接"`
	AppId      string       `json:"appId" gorm:"type:varchar(100);not null;default:'';comment:跳转到三方小程序时小程序appid"`
	Sort       int          `json:"sort" gorm:"type:int4;not null;default:0;comment:排序"`
	Redirect   string       `json:"redirect" gorm:"type:varchar(1000);not null;comment:跳转路径"`
	Status     BannerStatus `json:"status" gorm:"type:int2;not null;default:1;comment:状态 1上线 2下线"`
	Ext        string       `json:"ext" gorm:"type:varchar(1000);not null;default:'';comment:额外参数"`
	CreateTime model.Time   `json:"createTime" gorm:"type:timestamptz;not null;comment:创建时间"`
	UpdateTime model.Time   `json:"updateTime" gorm:"type:timestamptz;not null;comment:更新时间"`
}
type BannerStatus int8

const (
	BannerStatusOk   BannerStatus = 1 //上线
	BannerStatusDown BannerStatus = 2 //下线
)

type BannerType string

const (
	BannerTypeMini BannerType = "mini"
	BannerTypePath BannerType = "path"
	BannerTypeH5   BannerType = "h5"
)

type BannerScene string

const (
	BannerSceneHome    BannerScene = "home"
	BannerSceneEvent   BannerScene = "event"
	BannerSceneTopic   BannerScene = "topic"
	BannerSceneWelfare BannerScene = "welfare"
)

const OrderByBannerSortAsc OrderBy = "order_by_banner_sort_asc"

var (
	BannerStatusMap = map[BannerStatus]string{
		BannerStatusOk:   "上线",
		BannerStatusDown: "下线",
	}
	BannerSceneMap = map[BannerScene]string{
		BannerSceneHome:    "首页",
		BannerSceneEvent:   "携手",
		BannerSceneTopic:   "社区",
		BannerSceneWelfare: "公益",
	}
	BannerTypeMap = map[BannerType]string{
		BannerTypeMini: "第三方小程序",
		BannerTypePath: "内部小程序路径",
		BannerTypeH5:   "H5链接",
	}
)
