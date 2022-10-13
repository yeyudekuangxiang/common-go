package service

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mio/config"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"time"
)

type DuiBaActivityService struct {
	ctx  *context.MioContext
	repo *repository.DuiBaActivityRepository
}

func NewDuiBaActivityService(ctx *context.MioContext) *DuiBaActivityService {
	return &DuiBaActivityService{
		ctx:  ctx,
		repo: repository.NewDuiBaActivityRepository(ctx),
	}
}

func (srv DuiBaActivityService) FindActivity(activityId string) (*entity.DuiBaActivity, error) {
	activity := entity.DuiBaActivity{}
	err := app.DB.Where("activity_id = ?", activityId).First(&activity).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return &activity, nil
}

func (srv DuiBaActivityService) Create(dto srv_types.CreateDuiBaActivityDTO) error {
	//判断是否存在
	banner, err := srv.repo.GetExistOne(repotypes.GetDuiBaActivityExistDO{
		ActivityId: dto.ActivityId})
	if err != nil {
		return err
	}
	if banner.ID != 0 {
		return errors.New("activityId已存在")
	}
	bannerDo := entity.DuiBaActivity{
		Status:    entity.DuiBaActivityStatusYes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now()}
	if err := util.MapTo(dto, &bannerDo); err != nil {
		return err
	}
	return srv.repo.Create(&bannerDo)
}

func (srv DuiBaActivityService) Update(dto srv_types.UpdateDuiBaActivityDTO) error {
	//判断是否存在
	info, err := srv.repo.GetExistOne(repotypes.GetDuiBaActivityExistDO{
		Id: dto.Id})
	if err != nil {
		return err
	}
	if info.ID == 0 {
		return errors.New("activityId不存在")
	}
	//是否存在
	one, errInfo := srv.repo.GetExistOne(repotypes.GetDuiBaActivityExistDO{ActivityId: dto.ActivityId, NotId: dto.Id})
	if errInfo != nil {
		return errInfo
	}
	if one.ID != 0 {
		return errors.New("activityId已存在")
	}
	do := entity.DuiBaActivity{
		Status:    entity.DuiBaActivityStatusYes,
		UpdatedAt: time.Now()}
	if err := util.MapTo(dto, &do); err != nil {
		return err
	}
	return srv.repo.Save(&do)
}

func (srv DuiBaActivityService) GetPageList(dto srv_types.GetPageDuiBaActivityDTO) ([]entity.DuiBaActivity, int64, error) {
	bannerDo := repotypes.GetDuiBaActivityPageDO{
		Statue: dto.Status,
	}
	if err := util.MapTo(dto, &bannerDo); err != nil {
		return nil, 0, err
	}
	list, total, err := srv.repo.GetPageList(bannerDo)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (srv DuiBaActivityService) Delete(dto srv_types.DeleteDuiBaActivityDTO) error {
	//判断是否存在
	info, err := srv.repo.GetExistOne(repotypes.GetDuiBaActivityExistDO{
		Id: dto.Id})
	if err != nil {
		return err
	}
	if info.ID == 0 {
		return errors.New("activityId不存在")
	}
	do := repotypes.DeleteDuiBaActivityDO{
		UpdatedAt: time.Now(),
		Id:        dto.Id,
		Status:    entity.DuiBaActivityStatusNo,
	}
	return srv.repo.Delete(&do)
}

func (srv DuiBaActivityService) Show(dto srv_types.ShowDuiBaActivityDTO) (*api_types.DuiBaActivityShowVO, error) {
	//判断是否存在
	info, err := srv.repo.GetExistOne(repotypes.GetDuiBaActivityExistDO{
		Id: dto.Id})
	if err != nil {
		return nil, err
	}
	if info.ID == 0 {
		return nil, errors.New("activityId不存在")
	}
	//生成兑吧页面路径
	ActivityAppPath := srv.GetActivityAppPath(info.ActivityId, info.Cid, info.IsShare, info.IsPhone)
	//生成兑吧页面预览小程序码
	ActivityViewQrCode, _ := srv.GetActivityViewQrCode(ActivityAppPath)
	//获取兑吧h5免登录链接
	activityH5 := srv.GetActivityH5(info.ActivityId, info.Cid)
	//生成兑吧页面路径
	jumpAppH5 := srv.GetJumpAppH5(info.ActivityId, info.Cid, info.IsShare, info.IsPhone)

	return &api_types.DuiBaActivityShowVO{
		ID:            info.ID,
		NoLoginH5Link: activityH5,
		StaticH5Link:  jumpAppH5,
		InsideLink:    ActivityAppPath,
		EwmLink:       ActivityViewQrCode,
	}, nil
}

// GetActivityAppPath 生成兑吧页面路径
// activityId 活动id
// cid 渠道id
// needShare 页面是否可以分享 1可以分享 2不可以分享
// checkPhone 访问页面是否必须绑定手机号 1必须绑定 2不必须绑定
// 返回值 pages/duiba_v2/duiba/duiba-share/index?activityId=001&cid=12&bind=bind
func (srv DuiBaActivityService) GetActivityAppPath(activityId string, cid int64, needShare entity.DuiBaActivityIsShare, checkPhone entity.DuiBaActivityIsPhone) string {
	path := ""
	isCheckPhone := util.Ternary(checkPhone == entity.DuiBaActivityIsPhoneYes, "true", "false").String()
	if needShare == entity.DuiBaActivityIsShareYes {
		path = fmt.Sprintf("/pages/duiba_v2/duiba-share/index?activityId=%s&cid=%d&bind=%s", activityId, cid, isCheckPhone)
	} else {
		path = fmt.Sprintf("/pages/duiba_v2/duiba-not-share/index?activityId=%s&cid=%d&bind=%s", activityId, cid, isCheckPhone)
	}
	return path
}

// GetActivityViewQrCode 生成兑吧页面预览小程序码
// path 小程序路径 pages/duiba_v2/duiba/duiba-share/index?activityId=001&cid=12&bind=bind
// 返回时 小程序码字节数组
func (srv DuiBaActivityService) GetActivityViewQrCode(path string) ([]byte, error) {

	qrCode, err := NewQRCodeService().GetLimitedQRCodeRaw(path, 256)
	if err != nil {
		return nil, err
	}
	return qrCode, nil
}

// GetActivityH5 获取兑吧h5免登录链接
// activityId 获取url
// 返回值 h5免登陆链接 https://go-api.miotech.com/api/mp2c/duiba/h5?activityId=q8sd82besdsdsd
func (srv DuiBaActivityService) GetActivityH5(activityId string, cid int64) string {
	path := fmt.Sprintf("/api/mp2c/duiba/h5?activityId=%s&cid=%d", activityId, cid)
	return util.LinkJoin(config.Config.App.Domain, path)
}

// GetJumpAppH5 生成兑吧页面路径
// activityId 活动id
// cid 渠道id
// needShare 页面是否可以分享 1可以分享 2不可以分享
// checkPhone 访问页面是否必须绑定手机号 1必须绑定 2不必须绑定
// return https://cloud1-1g6slnxm1240a5fb-1306244665.tcloudbaseapp.com/duiba_share_v2.html?activityId=index&cid=12&bind=true
func (srv DuiBaActivityService) GetJumpAppH5(activityId string, cid int64, needShare entity.DuiBaActivityIsShare, checkPhone entity.DuiBaActivityIsPhone) string {
	link := ""
	checkPhoneParam := util.Ternary(checkPhone == entity.DuiBaActivityIsPhoneYes, "", "&bind=bind").String()
	if needShare == entity.DuiBaActivityIsShareYes {
		link = fmt.Sprintf("https://cloud1-1g6slnxm1240a5fb-1306244665.tcloudbaseapp.com/duiba_share_v2.html?activityId=%s&cid=%d%s", activityId, cid, checkPhoneParam)
	} else {
		link = fmt.Sprintf("https://cloud1-1g6slnxm1240a5fb-1306244665.tcloudbaseapp.com/duiba_not_share_v2.html?activityId=%s&cid=%d%s", activityId, cid, checkPhoneParam)
	}
	return link
}
