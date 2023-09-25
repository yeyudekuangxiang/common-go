package community

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"time"
)

// topic
type TopicChangeLikeResp struct {
	TopicTitle  string `json:"topic"`
	TopicId     int64  `json:"topicId"`
	TopicUserId int64  `json:"TopicUserId"`
	LikeStatus  int    `json:"likeStatus"`
	IsFirst     bool   `json:"isFirst"`
}

type TopicListParams struct {
	ID           int64              `json:"id"`
	Ids          []int64            `json:"ids"`
	Rids         []string           `json:"rids"`
	TopicTagId   int64              `json:"topicTagId"`
	Status       int                `json:"status"` //0全部 1待审核 2审核失败 3已发布 4已下架
	IsTop        int                `json:"isTop"`
	IsEssence    int                `json:"isEssence"`
	UserId       int64              `json:"userId"` // 用于查询用户对帖子是否点赞
	OrderByList  entity.OrderByList `json:"orderByList"`
	OrderBy      entity.OrderBy     `json:"orderBy"`
	Label        string             `json:"label"`
	Type         int                `json:"type"`
	ActivityType int                `json:"activityType"`
	Offset       int                `json:"offset"`
	Limit        int                `json:"limit"` // limit为0时不限制数量
}

type MyTopicListParams struct {
	UserId int64 `json:"userId"`
	Status int   `json:"status"`
	Type   int   `json:"type"`
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
}

type TopicDetailResp struct {
	entity.Topic
	IsLike        bool             `json:"isLike"`
	UpdatedAtDate string           `json:"updatedAtDate"` //03-01
	User          entity.ShortUser `json:"user"`
}

type CreateTopicParams struct {
	Title   string   `json:"title" `
	Content string   `json:"content"`
	Images  []string `json:"images" `
	TagIds  []int64  `json:"tagIds"`
	Type    int      `json:"type"`
	TopicActivityParams
}

type TopicActivityParams struct {
	Region         string  `json:"region"`
	Address        string  `json:"address" `
	SaTags         []saTag `json:"saTags"`
	Remarks        string  `json:"remarks"`
	Qrcode         string  `json:"qrcode"`
	MeetingLink    string  `json:"meetingLink"`
	Contacts       string  `json:"contacts"`
	StartTime      int64   `json:"startTime"`
	EndTime        int64   `json:"endTime"`
	SignupDeadline int64   `json:"signupDeadline"`
	ActivityType   int     `json:"activityType"`
	SignupNumber   int     `json:"signupNumber"` //报名数量上限
}

type saTag struct {
	Type     int      `json:"type"`
	Code     string   `json:"code"`
	Category int      `json:"category"`
	Title    string   `json:"title"`
	Options  []string `json:"options"`
}

type UpdateTopicParams struct {
	ID int64 `json:"id"`
	CreateTopicParams
}

type FindTopicParams struct {
	TopicId int64 `json:"topicId"`
	UserId  int64 `json:"userId"`
	Type    int   `json:"type,omitempty"`
	Status  int   `json:"status"`
}

type AdminTopicListParams struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	TagId         int64     `json:"tagId"`
	TagIds        []string  `json:"tagIds"`
	UserId        int64     `json:"userId"`
	UserName      string    `json:"userName"`
	Status        int       `json:"status"`
	IsTop         int       `json:"isTop"`
	IsEssence     int       `json:"isEssence"`
	IsPartners    int       `json:"isPartners"`
	Position      string    `json:"position"`
	Type          int       `json:"type"`
	ActivityType  int       `json:"activityType"`
	PushStartTime time.Time `json:"pushStartTime"`
	PushEndTime   time.Time `json:"pushEndTime"`
	Offset        int       `json:"offset"`
	Limit         int       `json:"limit"`
}

// activity
type SignupInfosParams struct {
	TopicId      int64        `json:"topicId"`
	UserId       int64        `json:"userId"`
	SignupInfos  []SignupInfo `json:"signupInfo"`
	SignupTime   time.Time    `json:"signupTime"`
	SignupStatus int          `json:"signupStatus"`
	//附加
	OpenId string `json:"openId"`
}

type SignupInfo struct {
	Title    string      `json:"title"`
	Code     string      `json:"code"`
	Category int         `json:"category"`
	Type     int         `json:"type"`
	Options  []string    `json:"options"`
	Value    interface{} `json:"value"`
}

// comment
func (a APIComment) ApiComment() *APICommentResp {
	return &APICommentResp{
		Id:        a.Id,
		Message:   a.Message,
		MemberId:  a.MemberId,
		Floor:     a.Floor,
		Count:     a.Count,
		RootCount: a.RootCount,
		LikeCount: a.LikeCount,
		HateCount: a.HateCount,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Member:    a.Member,
		IsAuthor:  a.IsAuthor,
	}
}

type APIComment struct {
	Id            string           `gorm:"primaryKey;autoIncrement" json:"id"`
	ObjId         string           `gorm:"type:int8;not null" json:"objId" json:"objId"` // 对象id （文章）
	ObjType       int8             `gorm:"type:int2" json:"objType"`                     // 保留字段
	Message       string           `gorm:"type:text;not null" json:"message"`
	MemberId      int64            `gorm:"type:int8;not null" json:"memberId"`                // 评论用户id
	RootCommentId int64            `gorm:"type:int8;not null:default:0" json:"rootCommentId"` // 根评论id，不为0表示是回复评论
	ToCommentId   int64            `gorm:"type:int8;not null:default:0" json:"toCommentId"`   // 父评论id，为0表示是根评论
	Floor         int32            `gorm:"type:int4;not null:default:0" json:"floor"`         // 评论楼层
	Count         int32            `gorm:"type:int4;not null:default:0" json:"count"`         // 该评论下评论总数
	RootCount     int32            `gorm:"type:int4;not null:default:0" json:"rootCount"`     // 该评论下根评论总数
	LikeCount     int32            `gorm:"type:int4;not null:default:0" json:"likeCount"`     // 该评论点赞总数
	HateCount     int32            `gorm:"type:int4;not null:default:0" json:"hateCount"`     // 该评论点踩总数
	State         int8             `gorm:"type:int2;not null:default:0" json:"state"`         // 状态 0-正常 1-隐藏
	DelReason     string           `gorm:"type:varchar(200)" json:"delReason"`                //删除理由
	Attrs         int8             `gorm:"type:int2" json:"attrs"`                            // 属性 00-正常 10-运营置顶 01-用户置顶 保留字段
	Version       int64            `gorm:"type:int8;version" json:"version"`                  // 版本号 保留字段
	CreatedAt     model.Time       `gorm:"createdAt" json:"createdAt"`
	UpdatedAt     model.Time       `gorm:"updatedAt" json:"updatedAt"`
	Member        entity.ShortUser `gorm:"foreignKey:ID;references:MemberId" json:"member"` // 评论用户
	IsAuthor      int8             `gorm:"type:int2" json:"isAuthor"`                       // 是否作者
	RootChild     []*APIComment    `gorm:"foreignKey:RootCommentId;association_foreignKey:Id" json:"rootChild"`
}

type APICommentResp struct {
	Id        string            `gorm:"primaryKey;autoIncrement" json:"id"`
	Message   string            `gorm:"type:text;not null" json:"message"`
	MemberId  int64             `gorm:"type:int8;not null" json:"memberId"`            // 评论用户id
	Floor     int32             `gorm:"type:int4;not null:default:0" json:"floor"`     // 评论楼层
	Count     int32             `gorm:"type:int4;not null:default:0" json:"count"`     // 该评论下评论总数
	RootCount int32             `gorm:"type:int4;not null:default:0" json:"rootCount"` // 该评论下根评论总数
	LikeCount int32             `gorm:"type:int4;not null:default:0" json:"likeCount"` // 该评论点赞总数
	HateCount int32             `gorm:"type:int4;not null:default:0" json:"hateCount"` // 该评论点踩总数
	CreatedAt model.Time        `json:"createdAt"`
	UpdatedAt model.Time        `json:"updatedAt"`
	Member    entity.ShortUser  `json:"member"`
	Detail    Detail            `json:"detail"`
	IsAuthor  int8              `json:"isAuthor"` // 是否作者
	IsLike    int               `json:"isLike"`
	RootChild []*APICommentResp `json:"rootChild"`
}

type Detail struct {
	ObjId       string `json:"objId"`
	ObjType     int64  `json:"objType"`
	ImageList   string `json:"imageList"`
	Description string `json:"description"`
}

type TurnCommentReq struct {
	UserId   int64  `json:"userId"`
	TurnType int    `json:"turnType"`
	TurnId   string `json:"turnId"`
}

type CommentCountResp struct {
	Date    time.Time
	TopicId int64
	Total   int64
}

type CommentChangeLikeResp struct {
	CommentMessage string `json:"topic"`
	CommentId      int64  `json:"topicId"`
	CommentUserId  int64  `json:"TopicUserId"`
	LikeStatus     int    `json:"likeStatus"`
	IsFirst        bool   `json:"isFirst"`
}

type FindListCountReq struct {
	TopicIds []int64 `json:"topicIds"`
}
