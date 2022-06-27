package service

import (
	"context"
	"errors"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repo_types"
	"mio/internal/pkg/service/service_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/timeutils"
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
		return "", errors.New("证书不存在")
	}
	switch cert.Type {
	case entity.CertTypeRandom:
		return util.RandomStr(badgeCodeLength, util.RandomStrUpper), nil
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
	currentDate := timeutils.Now().StartOfDay().String()

	code := strings.Builder{}
	code.WriteString(badgeCodePrefix)
	code.WriteString(currentDate)

	if len(badge.Code) >= 12 && badge.Code[4:12] == currentDate {
		lastNum, _ := strconv.Atoi(badge.Code[12:])
		if lastNum != 0 {
			lastNum += 1
			code.WriteString(strconv.Itoa(lastNum))
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
func (srv BadgeService) FindBadge(param service_types.FindBadgeParam) (*entity.Badge, error) {
	return srv.repo.FindBadge(repo_types.FindBadgeBy{
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
	badge, err := srv.repo.FindBadge(repo_types.FindBadgeBy{
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
	return srv.repo.GetBadgeList(repo_types.GetBadgeListBy{
		OpenId: openid,
	})
}
func (srv BadgeService) UpdateBadgeIsNew(openid string, id int64, isNew bool) error {
	badge, err := srv.repo.FindBadge(repo_types.FindBadgeBy{
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
