package activity

import (
	"errors"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/activity"
	activity2 "mio/internal/pkg/repository/activity"
)

type GetRecordPageListParam struct {
	UserId                  int64 `json:"userId"`
	ApplyStatus             int8  `json:"applyStatus"`
	ShareUserBocBonusStatus int8  `json:"shareUserBocBonusStatus"`
	Offset                  int   `json:"offset"`
	Limit                   int   `json:"limit"`
}
type AddApplyRecordParam struct {
	UserId      int64  `json:"userId"`
	ShareUserId int64  `json:"shareUserId"` //分享者用户id
	Source      string `json:"source"`
}

type BocRecordDetail struct {
	activity.BocRecord
	CreatedAtDate model.Date       `json:"createdAtDate"`
	UpdatedAtDate model.Date       `json:"updatedAtDate"`
	User          entity.ShortUser `json:"user"`
}

type CreateBocShareBonusRecordParam struct {
	UserId int64
	//金额 单位分
	Value int64
	Type  activity.BocShareBonusType
	Info  string
}
type AnswerGMQuestionParam struct {
	UserId  int64
	Title   string
	Answer  string
	IsRight bool
}

type GDDbHomePageResponse struct {
	User      activity2.GDDbHomePageUserInfo
	School    []activity.GDDbSchoolRank
	IsNewUser bool
}

//活动认证相关 start

type XingAccess struct {
	OperatorID        string `json:"operatorID,omitempty"`
	SuccStat          int    `json:"succStat,omitempty"`
	AccessToken       string `json:"accessToken,omitempty"`
	TokenAvailableTim int    `json:"tokenAvailableTim,omitempty"`
	FailReason        int    `json:"failReason,omitempty"`
}

type XingAccessRequest struct {
	OperatorID     string `json:"operatorID,omitempty"`
	OperatorSecret string `json:"operatorSecret,omitempty"`
}

//活动认证相关 end

var NotInvite error = errors.New("not invite")
