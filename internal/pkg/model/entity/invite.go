package entity

import "mio/internal/pkg/model"

type InviteType string

const (
	InviteTypeRegular    InviteType = "REGULAR"
	InviteTypeAdjustment InviteType = "ADJUSTMENT"
	InviteTypeGreenTorch InviteType = "GREEN_TORCH"
)

type Invite struct {
	ID              int64      `json:"id"`
	InvitedByOpenId string     `json:"invitedByOpenId" gorm:"column:invited_by_openid"`
	NewUserOpenId   string     `json:"newUserOpenId" gorm:"column:new_user_openid"`
	Time            model.Time `json:"time"`
	InviteType      InviteType `json:"inviteType"`
	InviteCode      string     `json:"inviteCode"`
}
