package community

import (
	"mio/internal/app/mp2c/controller"
	"time"
)

type topicCountRequest struct {
	TopicId int64 `json:"topicId"`
	Status  int   `json:"status"`
}

type ActivitiesTagPageRequest struct {
	controller.PageFrom
}

type ActivitiesTagListRequest struct {
}

type IdRequest struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}

type GetTopicPageListRequest struct {
	ID         int64  `json:"id" form:"id" binding:"gte=0" alias:"文章id"`
	TopicTagId int64  `json:"topicTagId" form:"topicTagId" binding:"gte=0" alias:"标签id"`
	Order      string `json:"order" form:"order" alias:"排序"`
	controller.PageFrom
}

type GetWeappQrCodeRequest struct {
	TopicId int64 `json:"topicId" form:"topicId" binding:"required" alias:"文章id"`
}

type ChangeTopicLikeRequest struct {
	TopicId int64 `json:"topicId" form:"topicId" binding:"required" alias:"文章id"`
}

type CreateTopicRequest struct {
	Title   string   `json:"title" form:"title" alias:"标题" binding:"required,min=2,max=64"`
	Content string   `json:"content" form:"content" alias:"内容" binding:"min=0,max=10000"`
	Images  []string `json:"images" form:"images" alias:"图片" binding:"required,min=1,max=12"`
	TagIds  []int64  `json:"tagIds" form:"tagIds" alias:"话题" binding:"max=2"`
	Type    int      `json:"type" form:"type" alias:"类型"`
	//报名活动字段
	TopicActivity
}

type TopicActivity struct {
	Region         string  `json:"region" form:"region" binding:"required_if=Type 1 ActivityType 1"`
	Address        string  `json:"address" form:"address" binding:"required_if=Type 1 ActivityType 1"`
	SATags         []saTag `json:"saTags" form:"saTags" binding:"required_if=Type 1"`
	Remarks        string  `json:"remarks" form:"remarks"`
	Qrcode         string  `json:"qrcode" form:"qrcode" binding:"required_if=ActivityType 2 MeetingLink ''"`
	MeetingLink    string  `json:"meetingLink" form:"meetingLink" binding:"required_if=ActivityType 2 Qrcode ''"`
	Contacts       string  `json:"contacts" form:"contacts" binding:"required_if=Type 1"`
	StartTime      int64   `json:"startTime" form:"startTime" binding:"required_if=Type 1"`
	EndTime        int64   `json:"endTime" form:"endTime" binding:"required_if=Type 1"`
	SignupDeadline int64   `json:"signupDeadline" form:"signupDeadline" binding:"required_if=Type 1"`
	ActivityType   int     `json:"activityType" form:"activityType" binding:"required_if=Type 1"`
	SignupNumber   int     `json:"signupNumber" form:"signupNumber" binding:"required_if=Type 1"`
}

type saTag struct {
	Type     int      `json:"type"`
	Code     string   `json:"code"`
	Category int      `json:"category"`
	Title    string   `json:"title"`
	Options  []string `json:"options"`
}

type UpdateTopicRequest struct {
	ID int64 `json:"id" form:"id" alias:"id" binding:"required,gte=1"`
	CreateTopicRequest
}

type MyTopicRequest struct {
	HomePageRequest
	Status int `json:"status" form:"status"`
	Type   int `json:"type" form:"type"`
	controller.PageFrom
}

type HomePageRequest struct {
	UserId int64 `json:"userId" form:"userId"`
}

// commond
type ListFormById struct {
	ID int64 `json:"id" form:"id" alias:"id" binding:"required,gte=1"`
	controller.PageFrom
}

// comment
type CommentCreateForm struct {
	Message string `json:"message" form:"message" alias:"评论内容" binding:"required,min=1"`
	Root    int64  `json:"root" form:"root" alias:"rootId" binding:"min=0"`
	Parent  int64  `json:"parent" form:"parent" alias:"parentId" binding:"min=0"`
	ObjId   int64  `json:"objId" form:"objId" alias:"objId" binding:"required,min=1"`
}

type CommentEditForm struct {
	CommentId int64  `json:"commentId" form:"commentId" alias:"commentId" binding:"required,min=1"`
	Message   string `json:"message" form:"message" alias:"评论内容" binding:"required,min=1"`
}

type ChangeCommentLikeForm struct {
	CommentId int64 `json:"commentId" form:"commentId" binding:"required" alias:"评论id"`
}

type TurnCommentRequest struct {
	TurnType int    `json:"turnType" form:"turnType" binding:"required"`
	TurnId   string `json:"turnId" form:"turnId" binding:"required"`
}

// collect
type MyCollectionRequest struct {
	controller.PageFrom
}

type CollectionRequest struct {
	ObjId   int64 `json:"objId" form:"objId" binding:"required"`
	ObjType int   `json:"objType" form:"objType"`
}

// activity
type SignupTopicRequest struct {
	TopicId     int64        `json:"topicId" form:"topicId" binding:"required"`
	SignupInfos []SignupInfo `json:"signupInfos" form:"signupInfos" binding:"required"`
}

// activity
type SignupTopicRequestV2 struct {
	TopicId  int64  `json:"topicId" form:"topicId" binding:"required"`
	RealName string `json:"realName" form:"realName"`
	Gender   int    `json:"gender" form:"gender"`
	Age      int    `json:"age" form:"age"`
	Phone    string `json:"phone" form:"phone" `
	Wechat   string `json:"wechat" form:"wechat"`
	City     string `json:"city" form:"city"`
	Remarks  string `json:"remarks" form:"remarks"`
}

type SignupInfo struct {
	Title    string      `json:"title"`
	Code     string      `json:"code"`
	Category int         `json:"category"`
	Type     int         `json:"type"`
	Options  []string    `json:"options"`
	Value    interface{} `json:"value"`
}

type MySignupRequest struct {
	controller.PageFrom
}

type SignupListResponse struct {
	SeeCount    int64    `json:"seeCount"`
	SignupCount int64    `json:"signupCount"`
	SignupList  []Signup `json:"signupList"`
}
type Signup struct {
	Id           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TopicId      int64     `json:"topicId"`
	UserId       int64     `json:"userId"`
	RealName     string    `json:"realName"`
	Phone        string    `json:"phone"`
	Gender       int       `json:"gender"`
	Age          int       `json:"age"`
	Wechat       string    `json:"wechat"`
	City         string    `json:"city"`
	Remarks      string    `json:"remarks"`
	SignupTime   time.Time `json:"signupTime"`
	CancelTime   time.Time `json:"cancelTime,omitempty"`
	SignupStatus int       `json:"signupStatus"`
}

type GetTagRequest struct {
	ID int64 `json:"id" form:"id" binding:"gte=0" alias:"tag id"`
	controller.PageFrom
}
