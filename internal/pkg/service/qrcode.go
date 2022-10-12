package service

import (
	"bytes"
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/skip2/go-qrcode"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
	"mio/pkg/errno"
	"mio/pkg/wxapp"
	"strings"
	"time"
)

type QRCodeService struct {
	repo *repository.QrCodeRepository
}

func NewQRCodeService() *QRCodeService {
	return &QRCodeService{repo: repository.NewQrCodeRepository(app.DB)}
}

// FindQrCode 查询二维码信息
// key  二维码的key值
// 返回值
// qr *entity.QRCode 二维码信息
// exist bool 二维码信息是否存在 true代表存在 false代表不存在
// err error 查询异常错误信息
func (srv QRCodeService) FindQrCode(key string) (qr *entity.QRCode, exist bool, err error) {
	return srv.repo.GetQrCode(key)
}

// CreateQrCode 将二维码信息存到数据库
func (srv QRCodeService) CreateQrCode(dto srv_types.CreateQrCodeDTO) (*entity.QRCode, error) {
	qr := entity.QRCode{}
	if err := util.MapTo(dto, &qr); err != nil {
		return nil, err
	}
	qr.QrCodeId = util.UUID()
	qr.CreatedAt = time.Now()
	if err := srv.repo.Create(&qr); err != nil {
		return nil, err
	}
	return &qr, nil
}

// GetUnlimitedQRCode 获取没有数量限制的小程序码
// qrScene entity.QrCodeScene 小程序码使用场景 必传
// page 小程序页面  pages/community/details/index 必传
// scene  小程序scene参数 id=11&age=12 不是必传
// width 小程序码宽度
// desc 小程序码描述 不是必传
// openId 记录生成者openid  不是必传
func (srv QRCodeService) GetUnlimitedQRCode(qrScene entity.QrCodeScene, page, scene string, width int, openId string) (*entity.QRCode, error) {
	key := srv.QrCodeKey(qrScene, page, scene)

	//判断是否已经创建过
	qr, exist, err := srv.FindQrCode(key)
	if err != nil {
		return nil, err
	}
	if exist {
		return qr, nil
	}

	//创建新的
	var resp *wxapp.QRCodeResponse
	err = app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		resp, err = app.Weapp.GetUnlimitedQRCodeResponse(&weapp.UnlimitedQRCode{
			Scene:     scene,
			Page:      page,
			Width:     width,
			IsHyaline: true,
		})
		if err != nil {
			return false, err
		}
		return app.Weapp.IsExpireAccessToken(resp.ErrCode)
	}, 1)

	if err != nil {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}
	if resp.ErrCode != 0 {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, resp)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	imagePath, err := DefaultOssService.PutObject(fmt.Sprintf("qrcode/%s/%s.png", qrScene, key), bytes.NewReader(resp.Buffer))
	if err != nil {
		//app.Logger.Errorf("上传分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	//保存数据
	qr, err = srv.CreateQrCode(srv_types.CreateQrCodeDTO{
		ImagePath:    imagePath,
		QrCodeScene:  qrScene,
		QrCodeSource: entity.QrCodeSourceWeappUnlimited,
		Key:          key,
		Content:      page,
		Ext:          scene,
		OpenId:       openId,
	})
	return qr, err
}

// GetUnlimitedQRCodeRaw 获取没有数量限制的小程序码字节数组
// page 小程序页面  pages/community/details/index 必传
// scene  小程序scene参数 id=11&age=12 不是必传
// width 小程序码宽度
func (srv QRCodeService) GetUnlimitedQRCodeRaw(page, scene string, width int) ([]byte, error) {
	//创建新的
	var resp *wxapp.QRCodeResponse
	err := app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		resp, err = app.Weapp.GetUnlimitedQRCodeResponse(&weapp.UnlimitedQRCode{
			Scene:     scene,
			Page:      page,
			Width:     width,
			IsHyaline: true,
		})
		if err != nil {
			return false, err
		}
		return app.Weapp.IsExpireAccessToken(resp.ErrCode)
	}, 1)

	if err != nil {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}
	if resp.ErrCode != 0 {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, resp)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	return resp.Buffer, nil
}

// GetLimitedQRCode 获取有数量限制的小程序码
// qrScene entity.QrCodeScene 小程序码使用场景 必传
// path 小程序路径  pages/community/details/index?id=11&age=12 必传
// width 小程序码宽度
// openId 记录生成者openid  不是必传
func (srv QRCodeService) GetLimitedQRCode(qrScene entity.QrCodeScene, path string, width int, openId string) (*entity.QRCode, error) {
	key := srv.QrCodeKey(qrScene, path)

	//判断是否已经创建过
	qr, exist, err := srv.FindQrCode(key)
	if err != nil {
		return nil, err
	}
	if exist {
		return qr, nil
	}

	//创建新的
	var resp *wxapp.QRCodeResponse
	err = app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		resp, err = app.Weapp.GetWxaCodeResponse(&weapp.QRCode{
			Path:      path,
			Width:     width,
			IsHyaline: true,
		})
		if err != nil {
			return false, err
		}
		return app.Weapp.IsExpireAccessToken(resp.ErrCode)
	}, 1)

	if err != nil {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}
	if resp.ErrCode != 0 {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, resp)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	imagePath, err := DefaultOssService.PutObject(fmt.Sprintf("qrcode/%s/%s.png", qrScene, key), bytes.NewReader(resp.Buffer))
	if err != nil {
		//app.Logger.Errorf("上传分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	//保存数据
	qr, err = srv.CreateQrCode(srv_types.CreateQrCodeDTO{
		ImagePath:    imagePath,
		QrCodeScene:  qrScene,
		QrCodeSource: entity.QrCodeSourceWeappLimited,
		Key:          key,
		Content:      path,
		OpenId:       openId,
	})
	return qr, err
}

// GetLimitedQRCodeRaw 获取有数量限制的小程序码字节数组
// path 小程序路径  pages/community/details/index?id=11&age=12 必传
// width 小程序码宽度
func (srv QRCodeService) GetLimitedQRCodeRaw(path string, width int) ([]byte, error) {

	//创建新的
	var resp *wxapp.QRCodeResponse
	err := app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		resp, err = app.Weapp.GetWxaCodeResponse(&weapp.QRCode{
			Path:      path,
			Width:     width,
			IsHyaline: true,
		})
		if err != nil {
			return false, err
		}
		return app.Weapp.IsExpireAccessToken(resp.ErrCode)

	}, 1)

	if err != nil {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}
	if resp.ErrCode != 0 {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, resp)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	return resp.Buffer, nil
}

// GetWxQrcode 获取有数量限制的小程序二维码
// qrScene entity.QrCodeScene 小程序码使用场景 必传
// path 小程序路径  pages/community/details/index?id=11&age=12 必传
// width 二维码宽度
// openId 记录生成者openid  不是必传
func (srv QRCodeService) GetWxQrcode(qrScene entity.QrCodeScene, path string, width int, openId string) (*entity.QRCode, error) {
	key := srv.QrCodeKey(qrScene, path)

	//判断是否已经创建过
	qr, exist, err := srv.FindQrCode(key)
	if err != nil {
		return nil, err
	}
	if exist {
		return qr, nil
	}

	//创建新的
	var resp *wxapp.QRCodeResponse
	err = app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		resp, err = app.Weapp.CreateWxaQrcodeResponse(&weapp.QRCodeCreator{
			Path:  path,
			Width: width,
		})
		if err != nil {
			return false, err
		}
		return app.Weapp.IsExpireAccessToken(resp.ErrCode)
	}, 1)

	if err != nil {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}
	if resp.ErrCode != 0 {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, resp)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	imagePath, err := DefaultOssService.PutObject(fmt.Sprintf("qrcode/%s/%s.png", qrScene, key), bytes.NewReader(resp.Buffer))
	if err != nil {
		//app.Logger.Errorf("上传分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	//保存数据
	qr, err = srv.CreateQrCode(srv_types.CreateQrCodeDTO{
		ImagePath:    imagePath,
		QrCodeScene:  qrScene,
		QrCodeSource: entity.QrCodeSourceWeappQr,
		Key:          key,
		Content:      path,
		OpenId:       openId,
	})
	return qr, err
}

// GetWxQrcodeRaw 获取有数量限制的小程序二维码字节数组
// path 小程序路径  pages/community/details/index?id=11&age=12 必传
// width 二维码宽度
func (srv QRCodeService) GetWxQrcodeRaw(path string, width int) ([]byte, error) {

	//创建新的
	var resp *wxapp.QRCodeResponse
	err := app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		resp, err = app.Weapp.CreateWxaQrcodeResponse(&weapp.QRCodeCreator{
			Path:  path,
			Width: width,
		})
		if err != nil {
			return false, err
		}
		return app.Weapp.IsExpireAccessToken(resp.ErrCode)
	}, 1)

	if err != nil {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}
	if resp.ErrCode != 0 {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, resp)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	return resp.Buffer, err
}

// GetTextQrCode 获取普通二维码
// qrScene entity.QrCodeScene 小程序码使用场景 必传
// content 二维码信息 https://baidu.com、hello world  必传
// width 二维码宽度
// openId 记录生成者openid  不是必传
func (srv QRCodeService) GetTextQrCode(qrScene entity.QrCodeScene, content string, width int, openId string) (*entity.QRCode, error) {
	key := srv.QrCodeKey(qrScene, content)

	//判断是否已经创建过
	qr, exist, err := srv.FindQrCode(key)
	if err != nil {
		return nil, err
	}
	if exist {
		return qr, nil
	}

	//创建新的
	qrData, err := qrcode.Encode(content, qrcode.Medium, width)

	if err != nil {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	imagePath, err := DefaultOssService.PutObject(fmt.Sprintf("qrcode/%s/%s.png", qrScene, key), bytes.NewReader(qrData))
	if err != nil {
		//app.Logger.Errorf("上传分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	//保存数据
	qr, err = srv.CreateQrCode(srv_types.CreateQrCodeDTO{
		ImagePath:    imagePath,
		QrCodeScene:  qrScene,
		QrCodeSource: entity.QrCodeSourceCommon,
		Key:          key,
		Content:      content,
		OpenId:       openId,
	})
	return qr, err

}

// GetTextQrCodeRaw 获取普通二维码字节数组
// content 二维码信息 https://baidu.com、hello world  必传
// width 二维码宽度
func (srv QRCodeService) GetTextQrCodeRaw(content string, width int) ([]byte, error) {
	//创建新的
	qrData, err := qrcode.Encode(content, qrcode.Medium, width)

	if err != nil {
		//app.Logger.Errorf("生成分享码失败 %v %v %+v\n", topicId, userId, err)
		return nil, errno.ErrCommon.WithMessage("生成分享码失败").WithErr(err)
	}

	return qrData, nil
}

func (srv QRCodeService) QrCodeKey(scene entity.QrCodeScene, content string, others ...string) string {
	keyStr := string(scene) + content + strings.Join(others, "")
	return encrypt.Md5(keyStr)
}
