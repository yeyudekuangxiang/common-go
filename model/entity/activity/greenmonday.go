package activity

import "mio/model"

type GMRecord struct {
	ID               int        `json:"id"`
	UserId           int64      `json:"userId"`
	AvailableQuesNum int        `json:"availableQuesNum"`
	UsedQuesNum      int        `json:"usedQuesNum"`
	RightQuesNum     int        `json:"rightQuesNUm"`
	WrongQuesNum     int        `json:"wrongQuesNum"`
	PrizeStatus      int        `json:"prizeStatus"` // 1没有获得领取资格 2未领取 3已领取
	CreatedAt        model.Time `json:"createdAt"`
	UpdatedAt        model.Time `json:"updatedAt"`
}

func (GMRecord) TableName() string {
	return "gm_record"
}

type GMInvitationRecord struct {
	ID               int        `json:"id"`
	UserId           int64      `json:"userId"`
	InviteeUserId    int64      `json:"inviteeUserId"`
	InviteeIsNewUser int        `json:"inviteeIsNewUser"` //被邀请人是否是新用户 1是新用户 2不是新用户
	CreatedAt        model.Time `json:"createdAt"`
	UpdatedAt        model.Time `json:"updatedAt"`
}

func (GMInvitationRecord) TableName() string {
	return "gm_invitation_record"
}

type GMQuestionLog struct {
	ID          int        `json:"id"`
	UserId      int64      `json:"userId"`
	Title       string     `json:"title"`
	Answer      string     `json:"answer"`
	IsRight     int        `json:"isRight"`     //1正确 2错误
	IsSendPoint int        `json:"isSendPoint"` //是否已发放答题积分 1未发放 2已发放
	CreatedAt   model.Time `json:"createdAt"`
	UpdatedAt   model.Time `json:"updatedAt"`
}

func (GMQuestionLog) TableName() string {
	return "gm_question_log"
}
