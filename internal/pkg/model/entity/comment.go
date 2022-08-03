package entity

import "time"

type CommentIndex struct {
	ID            int64  `gorm:"primaryKey;autoIncrement"`
	ObjId         int64  `gorm:"type:int8;not null"` // 对象id （文章）
	ObjType       int8   `gorm:"type:int2"`          // 保留字段
	Message       string `gorm:"type:text;not null"`
	MemberId      int64  `gorm:"type:int8;not null"`           // 评论用户id
	RootCommentId int64  `gorm:"type:int8;not null:default:0"` // 根评论id，不为0表示是回复评论
	ToCommentId   int64  `gorm:"type:int8;not null:default:0"` // 父评论id，为0表示是根评论
	ToNickName    string `gorm:"size:128;"`                    // 子评论的回复评论，回复目标的用户昵称
	Floor         int32  `gorm:"type:int4;not null:default:0"` // 评论楼层
	Count         int32  `gorm:"type:int4;not null:default:0"` // 该评论下评论总数
	RootCount     int32  `gorm:"type:int4;not null:default:0"` // 该评论下根评论总数
	LikeCount     int32  `gorm:"type:int4;not null:default:0"` // 该评论点赞总数
	HateCount     int32  `gorm:"type:int4;not null:default:0"` // 该评论点踩总数
	State         int8   `gorm:"type:int2;not null:default:0"` // 状态 0-正常 1-隐藏
	Attrs         int8   `gorm:"type:int2"`                    // 属性 00-正常 10-运营置顶 01-用户置顶 保留字段
	Version       int64  `gorm:"type:int8;version"`            // 版本号 保留字段
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Member        ShortUser      `gorm:"foreignKey:ID;references:MemberId"` // 评论用户
	RootChild     []CommentIndex `gorm:"foreignKey:RootCommentId;association_foreignKey:Id"`
}

func (CommentIndex) TableName() string {
	return "comment_index"
}
