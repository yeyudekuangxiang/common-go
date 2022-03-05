package entity

import "mio/model"

// 标签
type Tag struct {
	model.Model
	Name        string `gorm:"size:32;unique" json:"name" form:"name"`          // 名称
	Description string `gorm:"size:1024" json:"description" form:"description"` // 描述
	Logo        string `gorm:"size:1024" json:"logo" form:"logo"`               // 图标
	Sort        int    `gorm:"index:idx_sort_" json:"sort" form:"sort"`         // 排序编号
	Status      int    `gorm:"not null" json:"status" form:"status"`            // 状态
	CreateTime  int64  `json:"createTime" form:"createTime"`                    // 创建时间
}

type Topic struct {
	Id            int64  `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
	TopicTag      string `gorm:"type:int;not null:default:0" json:"topicTag" form:"topicTag"`     // 类型
	UserId        int64  `gorm:"not null;index:idx_topic_user_id;" json:"userId" form:"userId"`   // 用户
	Title         string `gorm:"size:128" json:"title" form:"title"`                              // 标题
	Content       string `gorm:"type:longtext" json:"content" form:"content"`                     // 内容
	ImageList     string `gorm:"type:longtext" json:"imageList" form:"imageList"`                 // 图片
	Recommend     bool   `gorm:"not null;index:idx_recommend" json:"recommend" form:"recommend"`  // 是否推荐
	RecommendTime int64  `gorm:"not null" json:"recommendTime" form:"recommendTime"`              // 推荐时间
	ViewCount     int64  `gorm:"not null" json:"viewCount" form:"viewCount"`                      // 查看数量
	CommentCount  int64  `gorm:"not null" json:"commentCount" form:"commentCount"`                // 跟帖数量
	LikeCount     int64  `gorm:"not null" json:"likeCount" form:"likeCount"`                      // 点赞数量
	Status        int    `gorm:"index:idx_topic_status;" json:"status" form:"status"`             // 状态：0：正常、1：删除
	CreateTime    int64  `gorm:"index:idx_topic_create_time" json:"createTime" form:"createTime"` // 创建时间
	Sort          int    `gorm:"index:idx_sort_" json:"sort" form:"sort"`                         // 排序编号
	Avatar        string `json:"avatar"`
	Tags          []Tag  `json:"tags" gorm:"many2many:topic_tag;"`
}

type TopicItemRes struct {
	model.Model
	TopicTag      string `gorm:"type:int;not null:default:0" json:"topicTag" form:"topicTag"`     // 类型
	UserId        int64  `gorm:"not null;index:idx_topic_user_id;" json:"userId" form:"userId"`   // 用户
	Title         string `gorm:"size:128" json:"title" form:"title"`                              // 标题
	Content       string `gorm:"type:longtext" json:"content" form:"content"`                     // 内容
	ImageList     string `gorm:"type:longtext" json:"imageList" form:"imageList"`                 // 图片
	Recommend     bool   `gorm:"not null;index:idx_recommend" json:"recommend" form:"recommend"`  // 是否推荐
	RecommendTime string `gorm:"not null" json:"recommendTime" form:"recommendTime"`              // 推荐时间
	ViewCount     int64  `gorm:"not null" json:"viewCount" form:"viewCount"`                      // 查看数量
	CommentCount  int64  `gorm:"not null" json:"commentCount" form:"commentCount"`                // 跟帖数量
	LikeCount     int64  `gorm:"not null" json:"likeCount" form:"likeCount"`                      // 点赞数量
	Status        int    `gorm:"index:idx_topic_status;" json:"status" form:"status"`             // 状态：0：正常、1：删除
	CreateTime    string `gorm:"index:idx_topic_create_time" json:"createTime" form:"createTime"` // 创建时间
}
