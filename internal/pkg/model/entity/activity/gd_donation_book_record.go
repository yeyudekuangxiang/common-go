package activity

import (
	"mio/internal/pkg/model"
)

//GDDonationBookRecord 参加活动人员状态及关系
type GDDonationBookRecord struct {
	ID             int64      `gorm:"primaryKey;not null;serial8;comment:id" json:"id"`
	UserId         int64      `gorm:"type:int8;not null:default:0" json:"userId"`
	AnswerStatus   int        `gorm:"type:int2;not null:default:0" json:"answerStatus"` //答题状态 0 未开始 1 进行中 2 已完成
	IsSuccess      int        `gorm:"type:int2;not null:default:0" json:"isSuccess"`    //是否成团 0 未成团 1 已成团
	InviteId       int64      `gorm:"type:int8;not null:default:0" json:"inviteId"`     //邀请人id
	InviteType     int        `gorm:"type:int2;not null:default:0" json:"inviteType"`   //邀请状态 0 团长 1 团员
	TitleUrl       string     `gorm:"type:varchar;not null:default:''" json:"titleUrl"`
	CertificateUrl string     `gorm:"type:varchar;not null:default:''" json:"certificateUrl"`
	CreatedAt      model.Time `json:"createdAt"`
	UpdatedAt      model.Time `json:"updatedAt"`
}

func (GDDonationBookRecord) TableName() string {
	return "gd_donation_book"
}

// GDDbUserSchool 参加活动人员学校信息 活动简称 gd：广东 db: donation book
type GDDbUserSchool struct {
	ID          int64      `gorm:"primaryKey;not null;serial8;comment:id" json:"id"`
	UserId      int64      `gorm:"type:int8;not null:default:0" json:"userId"`
	UserName    string     `gorm:"type:varchar(50);not null:default:''" json:"userName"`
	SchoolId    int64      `gorm:"type:int8;not null:default:0" json:"schoolId"`
	GradeId     int64      `gorm:"type:int8;not null:default:0" json:"gradeId"`     //年级id
	ClassNumber uint32     `gorm:"type:int8;not null:default:0" json:"classNumber"` //班级号码
	CreatedAt   model.Time `json:"createdAt"`
	UpdatedAt   model.Time `json:"updatedAt"`
}

func (GDDbUserSchool) TableName() string {
	return "gd_db_user_school"
}

//GDDbSchool 学校
type GDDbSchool struct {
	ID         int64      `gorm:"primaryKey;not null;serial8;comment:id" json:"id"`
	CityId     int64      `gorm:"type:int8;not null:default:0" json:"cityId"` //市id
	Type       int        `gorm:"type:int2;not null:default:0" json:"type"`   //0 未分类 1 小学 2 中学
	SchoolName string     `json:"schoolName"`
	CreatedAt  model.Time `json:"createdAt"`
	UpdatedAt  model.Time `json:"updatedAt"`
}

func (GDDbSchool) TableName() string {
	return "gd_db_school"
}

// GDDbSchoolRank 学校证书表
type GDDbSchoolRank struct {
	ID           int64      `gorm:"primaryKey;not null;serial8;comment:id" json:"id"`
	SchoolId     int64      `gorm:"type:int8;not null:default:0" json:"schoolId"`
	SchoolName   string     `gorm:"type:varchar(100);not null:default:''" json:"schoolName"`
	DonateNumber int64      `gorm:"type:int8;not null:default:0" json:"donateNumber"`
	CreatedAt    model.Time `json:"createdAt"`
	UpdatedAt    model.Time `json:"updatedAt"`
}

func (GDDbSchoolRank) TableName() string {
	return "gd_db_school_rank"
}

// GDDbCity 城市表
type GDDbCity struct {
	ID       int64  `gorm:"primaryKey;not null;serial8;comment:id" json:"id"`
	CityName string `gorm:"type:varchar(50);not null:default:''" json:"cityName"`
}

func (GDDbCity) TableName() string {
	return "gd_db_city"
}

//
type GDDbGrade struct {
	ID    int64  `gorm:"primaryKey;not null;serial8;comment:id" json:"id"`
	Grade string `gorm:"type:varchar(50);not null:default:''" json:"grade"`
	Type  int    `gorm:"type:int2;not null:default:0" json:"type"` //0 小学 1 中学
}

func (GDDbGrade) TableName() string {
	return "gd_db_grade"
}
