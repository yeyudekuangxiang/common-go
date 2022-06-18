package activity

import "mio/internal/pkg/model/entity/activity"

type FindRecordBy struct {
	UserId int64
}
type GetRecordListBy struct {
	ApplyStatus             int8 `binding:"oneof=0 1 2 3 4" alias:"申请状态"` //0 全部 1参与中 2申请中 3申请成功 4申请失败
	ShareUserBocBonusStatus int8 //0全部 1未领取 2已领取
	UserIds                 []int64
	ShareUserIds            []int64
	Offset                  int
	Limit                   int
}
type FindGMRecordBy struct {
	UserId int64
}
type FindGMQuesLogBy struct {
	UserId int64
}
type FindGMInvitationRecordBy struct {
	UserId        int64
	InviteeUserId int64
}

type GDDbHomePageUserInfo struct {
	UserInfo    GDDbUserInfo
	InviteInfo  GDDbUserInfo
	InvitedInfo []GDDbUserInfo
}

type GDDbUserInfo struct {
	activity.GDDonationBookRecord
	AvatarUrl string `json:"avatarUrl"`
	Nickname  string `json:"nickname"`
}

type GDDbUserSchool struct {
	GDDbUserInfo
	SchoolName  string `json:"schoolName"`
	UserName    string `json:"userName"`
	CityName    string `json:"cityName"`
	GradeName   string `json:"gradeName"`
	ClassNumber uint32 `json:"classNumber"`
}

type FindSchoolBy struct {
	SchoolId   int64   `json:"schoolId"`
	SchoolIds  []int64 `json:"schoolIds"`
	SchoolName string  `json:"schoolName"`
	GradeType  int     `json:"gradeType"`
	CityId     int64   `json:"cityId"`
}
