package entity

import (
	"mio/internal/pkg/model"
)

const (
	OrderByTopicSortDesc = "order_by_topic_sort_desc"
)

type TopicStatus int

const (
	TopicStatusNeedVerify   = 1 //待审核
	TopicStatusVerifyFailed = 2 //审核失败
	TopicStatusPublished    = 3 //已发布
	TopicStatusHidden       = 4 //已下架
)

type Topic struct {
	Id            int64          `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
	TopicTag      string         `gorm:"-;type:int;not null:default:0" json:"-" form:"topicTag"`        // 类型
	UserId        int64          `gorm:"not null;index:idx_topic_user_id;" json:"userId" form:"userId"` // 用户
	User          ShortUser      `json:"user"`
	Title         string         `gorm:"size:128" json:"title" form:"title"`                       // 标题
	Content       string         `gorm:"type:longtext" json:"content" form:"content"`              // 内容
	ImageList     string         `gorm:"type:longtext" json:"imageList" form:"imageList"`          // 图片
	Recommend     bool           `gorm:"-;not null;index:idx_recommend" json:"-" form:"recommend"` // 是否推荐
	RecommendTime int64          `gorm:"-;not null" json:"-" form:"recommendTime"`                 // 推荐时间
	ViewCount     int64          `gorm:"not null" json:"viewCount" form:"viewCount"`               // 查看数量
	CommentCount  int64          `gorm:"not null" json:"commentCount" form:"commentCount"`         // 跟帖数量
	LikeCount     int64          `gorm:"not null" json:"likeCount" form:"likeCount"`               // 点赞数量
	Status        TopicStatus    `gorm:"index:idx_topic_status;" json:"status" form:"status"`      // 状态：1 待审核 2审核失败 3已发布 4已下架
	Sort          int            `gorm:"index:idx_sort_" json:"sort" form:"sort"`                  // 排序编号
	Avatar        string         `gorm:"-" json:"-"`
	Nickname      string         `gorm:"-" json:"-"`
	Tags          []Tag          `json:"tags" gorm:"many2many:topic_tag;"`
	Comment       []CommentIndex `json:"comment" gorm:"foreignKey:ObjId"`
	TopicTagId    string         `gorm:"-" json:"-"` // 类型
	SeeCount      int            `json:"seeCount"`   //浏览次数
	CreatedAt     model.Time     `json:"createdAt"`
	UpdatedAt     model.Time     `json:"updatedAt"`
	ImportId      int            `json:"-"`
	IsTop         int            `json:"isTop"`     //是否置顶
	IsEssence     int            `json:"isEssence"` //是否精华
	DelReason     string         `json:"delReason"` //审核不通过 or 删除 的理由
	IsLike        int            `json:"isLike"`    //是否点赞
}

func (Topic) TableName() string {
	return "topic"
}

type TopicItemRes struct {
	Id           int64          `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
	UserId       int64          `gorm:"not null;index:idx_topic_user_id;" json:"userId" form:"userId"` // 用户
	Title        string         `gorm:"size:128" json:"title" form:"title"`                            // 标题
	Content      string         `gorm:"type:longtext" json:"content" form:"content"`                   // 内容
	ImageList    string         `gorm:"type:longtext" json:"imageList" form:"imageList"`               // 图片
	ViewCount    int64          `gorm:"not null" json:"viewCount" form:"viewCount"`                    // 查看数量
	CommentCount int64          `gorm:"not null" json:"commentCount" form:"commentCount"`              // 跟帖数量
	LikeCount    int64          `gorm:"not null" json:"likeCount" form:"likeCount"`                    // 点赞数量
	Status       int            `gorm:"index:idx_topic_status;" json:"status" form:"status"`           // 状态：0正常 1禁言 2封号
	User         ShortUser      `json:"user"`
	Tags         []Tag          `json:"tags" gorm:"many2many:topic_tag;"`
	Comment      []CommentIndex `json:"comment" gorm:"foreignKey:ObjId"`
	IsTop        int            `json:"isTop"`     //是否置顶
	IsEssence    int            `json:"isEssence"` //是否精华
	DelReason    string         `json:"delReason"` //审核不通过 or 删除 的理由
}

func (TopicItemRes) TableName() string {
	return "topic"
}
