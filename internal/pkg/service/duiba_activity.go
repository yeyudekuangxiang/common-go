package service

import (
	"fmt"
	"gorm.io/gorm"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util"
)

var DefaultDuiBaActivityService = DuiBaActivityService{}

type DuiBaActivityService struct {
}

func (srv DuiBaActivityService) FindActivity(activityId string) (*entity.DuiBaActivity, error) {
	activity := entity.DuiBaActivity{}
	err := app.DB.Where("activity_id = ?", activityId).First(&activity).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	return &activity, nil
}

// GetActivityAppPath 生成兑吧页面路径
// activityId 活动id
// cid 渠道id
// needShare 页面是否可以分享 1可以分享 2不可以分享
// checkPhone 访问页面是否必须绑定手机号 1必须绑定 2不必须绑定
// 返回值 pages/duiba_v2/duiba/duiba-share/index?activityId=001&cid=12&bind=bind
func (srv DuiBaActivityService) GetActivityAppPath(activityId string, cid string, needShare, checkPhone int) string {
	path := ""
	checkPhoneParam := util.Ternary(checkPhone == 1, "", "&bind=bind")
	if needShare == 1 {
		path = fmt.Sprintf("pages/duiba_v2/duiba/duiba-share/index?activityId=%s&cid=%s%s", activityId, cid, checkPhoneParam)
	} else {
		path = fmt.Sprintf("pages/duiba_v2/duiba-not-share/index?activityId=%s&cid=%s%s", activityId, cid, checkPhoneParam)
	}
	return path
}

// GetActivityViewQrCode 生成兑吧页面预览小程序码
// path 小程序路径 pages/duiba_v2/duiba/duiba-share/index?activityId=001&cid=12&bind=bind
// 返回时 https://resources.miotech.com/qrcode/activity_view/sdoiojadsiod.png
func (srv DuiBaActivityService) GetActivityViewQrCode(path string) (string, error) {

	qrCode, err := NewQRCodeService().GetLimitedQRCode(entity.QrCodeSceneDuibaView, path, 256, "")
	if err != nil {
		return "", err
	}
	return DefaultOssService.FullUrl(qrCode.ImagePath), nil
}

// GetActivityH5 获取兑吧h5免登录链接
// activityId 获取url
// 返回值 h5免登陆链接 https://go-api.miotech.com/api/mp2c/duiba/h5?activityId=q8sd82besdsdsd
func (srv DuiBaActivityService) GetActivityH5(activityId string) string {
	return util.LinkJoin(config.Config.App.Domain, "/api/mp2c/duiba/h5?activityId="+activityId)
}

// GetJumpAppH5 生成兑吧页面路径
// activityId 活动id
// cid 渠道id
// needShare 页面是否可以分享 1可以分享 2不可以分享
// checkPhone 访问页面是否必须绑定手机号 1必须绑定 2不必须绑定
// return https://cloud1-1g6slnxm1240a5fb-1306244665.tcloudbaseapp.com/duiba_share_v2.html?activityId=index&cid=12&bind=true
func (srv DuiBaActivityService) GetJumpAppH5(activityId string, cid string, needShare, checkPhone int) string {
	link := ""
	checkPhoneParam := util.Ternary(checkPhone == 1, "", "&bind=bind")
	if needShare == 1 {
		link = fmt.Sprintf("https://cloud1-1g6slnxm1240a5fb-1306244665.tcloudbaseapp.com/duiba_share_v2.html?activityId=%s&cid=%s%s", activityId, cid, checkPhoneParam)
	} else {
		link = fmt.Sprintf("https://cloud1-1g6slnxm1240a5fb-1306244665.tcloudbaseapp.com/duiba_not_share_v2.html?activityId=%s&cid=%s%s", activityId, cid, checkPhoneParam)
	}
	return link
}
