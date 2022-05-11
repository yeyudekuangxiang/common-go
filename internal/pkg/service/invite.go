package service

import (
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util"
)

var DefaultInviteService = InviteService{}

type InviteService struct {
}

func (srv InviteService) GetInviteQrCode(openid string) (*QrCodeInfo, error) {
	qrcode := entity.QRCode{}
	err := app.DB.Where("openid = ? and qr_code_type = ?", openid, entity.QrCodeTypeSHARE).Order("id desc").Take(&qrcode).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	if qrcode.ID != 0 {
		return &QrCodeInfo{
			QrCodeId:    qrcode.QrCodeId,
			OpenId:      qrcode.OpenId,
			ImageUrl:    qrcode.ImageUrl,
			Description: qrcode.Description,
		}, nil
	}

	resp, comErr, err := app.Weapp.GetQRCode(&weapp.QRCode{
		Path: "/pages/home/index?invitedBy=" + openid,
	})
	if err != nil {
		return nil, err
	}
	if comErr.ErrCode != 0 {
		return nil, errors.New(comErr.ErrMSG)
	}
	defer resp.Body.Close()
	imageUrl, err := DefaultOssService.PutObject(fmt.Sprintf("mp2c/qrcode/share/%s.png", util.UUID()), resp.Body)
	if err != nil {
		return nil, err
	}
	qrcode = entity.QRCode{
		QrCodeId:    util.UUID(),
		ImageUrl:    imageUrl,
		QrCodeType:  entity.QrCodeTypeSHARE,
		OpenId:      openid,
		Description: "Share to friends",
	}
	app.DB.Create(&qrcode)

	return &QrCodeInfo{
		QrCodeId:    qrcode.QrCodeId,
		OpenId:      qrcode.OpenId,
		ImageUrl:    qrcode.ImageUrl,
		Description: qrcode.Description,
	}, nil
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
			Time:      user.Time.Date(),
		})
	}
	return infoList, nil
}
func (srv InviteService) AddInvite(openid, InvitedByOpenId string) (*entity.Invite, bool, error) {
	invite := entity.Invite{}
	err := app.DB.Where("new_user_openid = ? and invited_by_openid <> ''", openid).First(&invite).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	if invite.ID != 0 {
		return &invite, false, nil
	}

	invite = entity.Invite{
		InvitedByOpenId: InvitedByOpenId,
		NewUserOpenId:   openid,
		Time:            model.NewTime(),
		InviteType:      entity.InviteTypeRegular,
		InviteCode:      "",
	}
	return &invite, true, app.DB.Create(&invite).Error
}
