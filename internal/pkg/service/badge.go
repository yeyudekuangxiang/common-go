package service

import (
	"errors"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/timeutils"
	"strconv"
	"strings"
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
