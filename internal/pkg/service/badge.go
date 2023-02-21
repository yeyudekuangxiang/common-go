package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/randomtool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/event"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"strconv"
	"strings"
	"time"
)

const badgeCodeLength = 10
const badgeCodePrefix = "WBCF"

var DefaultBadgeService = BadgeService{repo: repository.DefaultBadgeRepository}

type BadgeService struct {
	repo repository.BadgeRepository
}

func (srv BadgeService) GenerateBadge(param GenerateBadgeParam) (*entity.Badge, error) {
	code, err := srv.GenerateCode(param.CertificateId)
	if err != nil {
		return nil, err
	}
	badge := entity.Badge{
		Code:          code,
		OpenId:        param.OpenId,
		CertificateId: param.CertificateId,
		ProductItemId: param.ProductItemId,
		CreateTime:    model.NewTime(),
		Partnership:   param.Partnership,
		OrderId:       param.OrderId,
	}
	return &badge, srv.repo.Create(&badge)
}
func (srv BadgeService) GenerateCode(certificateId string) (string, error) {
	cert, err := DefaultCertificateService.FindCertificate(FindCertificateBy{
		CertificateId: certificateId,
	})
	if err != nil {
		return "", err
	}
	if cert.ID == 0 {
		return "", errno.ErrCommon.WithMessage("证书不存在")
	}
	switch cert.Type {
	case entity.CertTypeRandom:
		return randomtool.RandomStr(badgeCodeLength, randomtool.RandomStrUpper), nil
	case entity.CertTypeRule:
		return srv.GenerateRuleCode(), nil
	case entity.CertTypeStock:
		stock, err := DefaultCertificateStockService.FindUnusedCertificate(cert.CertificateId)
		if err != nil {
			return "", err
		}
		return stock.Code, nil
	}
	return "", nil
}
func (srv BadgeService) GenerateRuleCode() string {
	badge := srv.repo.FindLastWithType(entity.CertTypeRule)
	currentDate := timetool.NowDate().Format("20060102")

	code := strings.Builder{}
	code.WriteString(badgeCodePrefix)
	code.WriteString(currentDate)

	/*WBCF202201250002*/
	if len(badge.Code) >= 12 && badge.Code[4:12] == currentDate {
		lastNum, _ := strconv.Atoi(badge.Code[12:])
		if lastNum != 0 {
			lastNum += 1
			code.WriteString(fmt.Sprintf("%04d", lastNum))
			return code.String()
		}
	}
	code.WriteString("0001")
	return code.String()
}
func (srv BadgeService) GetUserCertCount(openId string) (int64, error) {
	return srv.repo.FindUserCertCount(openId)
}
func (srv BadgeService) GetUserCertCountById(userId int64) (int64, error) {
	user, err := DefaultUserService.GetUserById(userId)
	if err != nil {
		return 0, err
	}
	if user.ID == 0 {
		return 0, nil
	}
	return srv.GetUserCertCount(user.OpenId)
}
func (srv BadgeService) FindBadge(param srv_types.FindBadgeParam) (*entity.Badge, error) {
	return srv.repo.FindBadge(repotypes.FindBadgeBy{
		OrderId: param.OrderId,
	})
}
func (srv BadgeService) UpdateCertImage(openid string, code string, imageUrl string) error {
	if !util.DefaultLock.Lock("UpdateCertImage:"+openid, time.Second*10) {
		return errno.ErrLimit.WithCaller()
	}
	defer util.DefaultLock.UnLock("UpdateCertImage:" + openid)

	badgeId, err := app.Redis.Get(context.Background(), config.RedisKey.BadgeImageCode+code).Int64()
	if err != nil {
		app.Logger.Error(err)
	}
	if badgeId == 0 {
		return errno.ErrTimeout.WithCaller()
	}
	badge, err := srv.repo.FindBadge(repotypes.FindBadgeBy{
		ID: badgeId,
	})
	if err != nil {
		return err
	}

	if badge.ID == 0 {
		return errno.ErrRecordNotFound
	}
	if badge.OpenId != openid {
		return errno.ErrInternalServer
	}
	badge.ImageUrl = imageUrl

	app.Redis.Del(context.Background(), config.RedisKey.BadgeImageCode+code)

	return srv.repo.Save(badge)
}
func (srv BadgeService) GetBadgePageList(openid string) ([]entity.Badge, error) {
	return srv.repo.GetBadgeList(repotypes.GetBadgeListBy{
		OpenId: openid,
	})
}
func (srv BadgeService) UpdateBadgeIsNew(openid string, id int64, isNew bool) error {
	badge, err := srv.repo.FindBadge(repotypes.FindBadgeBy{
		ID: id,
	})
	if err != nil {
		return err
	}
	if badge.ID == 0 {
		return errno.ErrRecordNotFound.WithCaller()
	}
	if badge.OpenId != openid {
		return errno.ErrInternalServer
	}

	badge.IsNew = isNew

	return srv.repo.Save(badge)
}

func (srv BadgeService) GetUploadOldBadgeSetting(userId int64, badgeId int64) (*srv_types.UploadBadgeResult, error) {
	userInfo, err := DefaultUserService.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	badge, err := srv.repo.FindBadge(repotypes.FindBadgeBy{
		ID:     badgeId,
		OpenId: userInfo.OpenId,
	})

	if err != nil {
		return nil, err
	}
	if badge.ID == 0 || badge.ProductItemId == "" {
		return nil, errno.ErrCommon.WithMessage("获取证书信息失败,请稍后再试")
	}
	if badge.ImageUrl != "" {
		return nil, errno.ErrCommon.WithMessage("证书已存在")
	}

	ev, err := event.DefaultEventService.FindEvent(event.FindEventParam{
		ProductItemId: badge.ProductItemId,
	})
	if err != nil {
		return nil, err
	}

	if ev.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("获取证书信息失败,请稍后再试")
	}

	setting, err := event.DefaultEventTemplateService.ParseSetting(ev.EventTemplateType, ev.TemplateSetting)
	if err != nil {
		app.Logger.Error(ev.EventTemplateType, ev.TemplateSetting, err)
		return nil, errors.New("系统异常,请稍后再试")
	}

	code := util.UUID()
	app.Redis.Set(context.Background(), config.RedisKey.BadgeImageCode+code, badge.ID, time.Minute*5)
	return &srv_types.UploadBadgeResult{
		EventTemplateType: ev.EventTemplateType,
		TemplateSetting: map[string]interface{}{
			string(ev.EventTemplateType): setting,
		},
		BadgeInfo: srv_types.UploadBadgeInfo{
			UploadCode:    code,
			Username:      userInfo.Nickname,
			Time:          badge.CreateTime.Format("2006-01-02 15:04:05"),
			CertificateNo: badge.Code,
		},
	}, nil

}
