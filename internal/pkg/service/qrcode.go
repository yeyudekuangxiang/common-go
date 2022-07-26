package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"mio/pkg/wxapp"
	"time"
)

type QRCodeService struct {
	repo *repository.QrCodeRepository
}

func NewQRCodeService() *QRCodeService {
	return &QRCodeService{repo: repository.NewQrCodeRepository(app.DB)}
}

//GetQrCode 获取二维码
// scene entity.QrCodeScene 二维码的使用场景
// key  二维码的key值 key值和scene应该组成唯一索引
// 返回值
// qr *entity.QRCode 二维码信息
// exist bool 二维码信息是否存在 true代表存在 false代表不存在
// err error 查询异常错误信息
func (srv QRCodeService) GetQrCode(scene entity.QrCodeScene, key string) (qr *entity.QRCode, exist bool, err error) {
	return srv.repo.GetQrCode(scene, key)
}
func (srv QRCodeService) CreateQrCode(dto srv_types.CreateQrCodeDTO) (*entity.QRCode, error) {
	qrcode := entity.QRCode{}
	if err := util.MapTo(dto, &qrcode); err != nil {
		return nil, err
	}
	qrcode.QrCodeId = util.UUID()
	qrcode.CreatedAt = time.Now()
	if err := srv.repo.Create(&qrcode); err != nil {
		return nil, err
	}
	return &qrcode, nil
}

// CreateTopicShareQr 创建酷喵圈分享小程序码
func (srv QRCodeService) CreateTopicShareQr(topicId int64, userId int64) (*entity.QRCode, error) {
	userInfo, err := DefaultUserService.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	page := "pages/cool-mio/mio-detail/index"
	scene := fmt.Sprintf("tid=%d&uid=%d&s=p", topicId, userId)

	key := util.Md5(userInfo.OpenId + scene)

	//判断是否已经创建过
	qr, exist, err := srv.GetQrCode(entity.QrCodeSceneTopicShare, key)
	if err != nil {
		return nil, err
	}
	if exist {
		return qr, nil
	}

	//创建新的
	resp, err := wxapp.NewClient(app.Weapp).GetUnlimitedQRCodeResponse(&weapp.UnlimitedQRCode{
		Scene:     scene,
		Page:      page,
		Width:     100,
		IsHyaline: true,
	})

	if err != nil {
		app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败,请稍后再试")
	}
	if resp.ErrCode != 0 {
		app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, resp)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败,请稍后再试")
	}

	path, err := DefaultOssService.PutObject(fmt.Sprintf("qrcode/topic-share/%s.png", key), bytes.NewReader(resp.Buffer))
	if err != nil {
		app.Logger.Errorf("上传分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败,请稍后再试")
	}

	//保存数据
	qr, err = srv.CreateQrCode(srv_types.CreateQrCodeDTO{
		ImagePath:    path,
		QrCodeScene:  entity.QrCodeSceneTopicShare,
		QrCodeSource: entity.QrCodeSourceWeappUnlimited,
		Key:          key,
		Content:      page,
		Ext:          scene,
		OpenId:       userInfo.OpenId,
		Description:  "酷喵圈文章分享海报小程序码",
	})
	return qr, err
}

// CreateInvite 创建邀请的积分小程序码
func (srv QRCodeService) CreateInvite(openId string) (*entity.QRCode, error) {

	page := "/pages/home/index?invitedBy=" + openId

	key := util.Md5(openId)

	//判断是否已经创建过
	qr, exist, err := srv.GetQrCode(entity.QrCodeSceneInvite, key)
	if err != nil {
		return nil, err
	}
	if exist {
		return qr, nil
	}

	resp, comErr, err := app.Weapp.GetQRCode(&weapp.QRCode{
		Path: page,
	})

	if err != nil {
		return nil, err
	}
	if comErr.ErrCode != 0 {
		return nil, errors.New(comErr.ErrMSG)
	}

	defer resp.Body.Close()
	imgPath, err := DefaultOssService.PutObject(fmt.Sprintf("qrcode/invite/%s.png", key), resp.Body)

	if err != nil {
		app.Logger.Errorf("生成分享码失败 %v %v\n", openId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败,请稍后再试")
	}

	//保存数据
	qr, err = srv.CreateQrCode(srv_types.CreateQrCodeDTO{
		ImagePath:    imgPath,
		QrCodeScene:  entity.QrCodeSceneInvite,
		QrCodeSource: entity.QrCodeSourceWeappLimited,
		Key:          key,
		Content:      page,
		OpenId:       openId,
		Description:  "邀请得积分小程序码",
	})
	return qr, err
}
