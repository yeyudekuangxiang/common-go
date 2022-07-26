package service

import (
	"fmt"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
)

var DefaultInviteService = InviteService{}

type InviteService struct {
}

func (srv InviteService) GetInviteList(openid string) ([]InviteInfo, error) {
	inviteList := make([]entity.Invite, 0)
	err := app.DB.Where("invited_by_openid = ? and time > '2022-01-10 00:00:00'", openid).Order("time desc").Find(&inviteList).Error
	if err != nil {
		panic(err)
	}
	infoList := make([]InviteInfo, 0)
	for _, invite := range inviteList {
		user, err := DefaultUserService.GetUserByOpenId(invite.NewUserOpenId)
		if err != nil {
			return nil, err
		}

		infoList = append(infoList, InviteInfo{
			OpenId:    user.OpenId,
			Nickname:  user.Nickname,
			AvatarUrl: user.AvatarUrl,
			Time:      user.Time,
		})
	}
	return infoList, nil
}
func (srv InviteService) AddInvite(openid, invitedByOpenId string) (*entity.Invite, bool, error) {
	if invitedByOpenId == "" {
		return &entity.Invite{}, false, nil
	}

	invite := entity.Invite{}
	err := app.DB.Where("new_user_openid = ? and invited_by_openid <> ''", openid).First(&invite).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	if invite.ID != 0 {
		return &invite, false, nil
	}

	invite = entity.Invite{
		InvitedByOpenId: invitedByOpenId,
		NewUserOpenId:   openid,
		Time:            model.NewTime(),
		InviteType:      entity.InviteTypeRegular,
		InviteCode:      "",
	}
	return &invite, true, app.DB.Create(&invite).Error
}
func (srv InviteService) Invite(openid, InvitedByOpenId string) error {
	app.Logger.Info("添加邀请关系", openid, InvitedByOpenId)
	_, isNew, err := srv.AddInvite(openid, InvitedByOpenId)
	if err != nil {
		return err
	}
	if !isNew {
		return nil
	}
	app.Logger.Info("发放邀请积分", openid, InvitedByOpenId, entity.PointCollectValueMap[entity.POINT_INVITE])
	//发放积分奖励
	_, err = DefaultPointTransactionService.Create(CreatePointTransactionParam{
		OpenId:       InvitedByOpenId,
		Type:         entity.POINT_INVITE,
		Value:        entity.PointCollectValueMap[entity.POINT_INVITE],
		AdditionInfo: fmt.Sprintf("invite %s", openid),
	})
	if err != nil {
		app.Logger.Error("发放邀请积分失败", err)
	}
	return err
}
