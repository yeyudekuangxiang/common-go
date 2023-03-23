package service

import (
	"context"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/baidu"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"mio/config"
	"mio/internal/pkg/core/app"
	mioctx "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	repo "mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/factory"
	"mio/pkg/errno"
	"time"
)

type OCRService struct {
	ctx         *mioctx.MioContext
	imageClient *baidu.ImageClient
	scanRepo    *repo.ScanLogRepository
}

func NewDefaultImageClient() *baidu.ImageClient {
	return factory.NewBaiDuImageFromTokenCenterRpc("baiduimage", app.RpcService.TokenCenterRpcSrv)
}

func NewOCRService(mioContext *mioctx.MioContext, imageClient *baidu.ImageClient) *OCRService {
	return &OCRService{
		ctx:         mioContext,
		imageClient: imageClient,
		scanRepo:    repo.NewScanLogRepository(mioContext),
	}
}
func DefaultOCRService() *OCRService {
	return NewOCRService(mioctx.NewMioContext(), NewDefaultImageClient())
}
func (srv OCRService) CheckIdempotent(openId string) error {
	if !util.DefaultLock.Lock("OCRIdempotent:"+openId, time.Second*10) {
		return errno.ErrLimit
	}
	return nil
}

func (srv OCRService) CheckRisk(risk int) error {
	if risk > 2 {
		return errno.ErrCommon.WithMessage("风险等级检测异常，请您稍后再试")
	}
	return nil
}

// OCRForGm 素食打卡
func (srv OCRService) OCRForGm(openid string, risk int, src string) error {
	err := srv.CheckIdempotent(openid)
	if err != nil {
		return err
	}
	err = srv.CheckRisk(risk)
	if err != nil {
		return err
	}
	res := util.OCRPush(src)
	var orderNo, fee string

	for k, v := range res.WordsResult {
		if v.Words == "收款:" {
			fee = res.WordsResult[k+1].Words
		}
		if v.Words == "联系电话:021-62333696" {
			orderNo = res.WordsResult[k+1].Words
		}
	}
	if orderNo == "" || fee == "" {
		return errno.ErrCommon.WithMessage("无法识别此账单,请重新上传,谢谢")
	}
	fmt.Println("素食打卡账单:" + orderNo + ":" + fee)
	cmd := app.Redis.SetNX(context.Background(), config.RedisKey.Lock+orderNo, "a", 36500*time.Second)
	if !cmd.Val() {
		fmt.Println(config.RedisKey.Lock + orderNo + "重复扫描素食打卡")
		return errno.ErrCommon.WithMessage("重复扫描素食打卡账单")
	}

	pointTranService := NewPointService(mioctx.NewMioContext())
	//发放积分
	_, err = pointTranService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       openid,
		ChangePoint:  100,
		BizId:        util.UUID(),
		Type:         entity.POINT_ADJUSTMENT,
		AdditionInfo: `{"素食打卡":"` + orderNo + `"}`,
	})

	return err
}

// Scan 扫描图片  此方法不会更新扫描次数 想要更新扫描次数 请使用 ScanWithHash 方法
func (srv OCRService) Scan(imgUrl string) ([]string, error) {
	rest, err := srv.imageClient.WebImage(baidu.WebImageParam{
		ImageUrl: imgUrl,
	})
	if err != nil {
		return nil, err
	}
	if !rest.IsSuccess() {
		return nil, errno.ErrCommon.WithMessage(rest.ErrorMsg)
	}

	results := make([]string, 0)
	for _, word := range rest.WordsResult {
		results = append(results, word.Words)
	}
	return results, nil
}

// ScanWithHash 扫描图片 并且更新扫描次数
func (srv OCRService) ScanWithHash(imageUrl string, imageHash string) ([]string, error) {
	result, err := srv.Scan(imageUrl)
	if err != nil {
		return nil, err
	}
	_, err = srv.UpdateImageScanCount(imageHash, imageUrl, result)
	if err != nil {
		app.Logger.Error("更新ocr扫描次数失败", imageUrl, imageHash, err)
	}
	return result, nil
}
func (srv OCRService) GetImageHash(imageUrl string) (string, error) {
	data, err := httptool.Get(imageUrl)
	if err != nil {
		return "", err
	}
	return encrypttool.Sha256Byte(data), nil
}
func (srv OCRService) UpdateImageScanCount(imageHash string, imageUrl string, scanResult []string) (int, error) {
	log, exist, err := srv.scanRepo.FindByHash(imageHash)
	if err != nil {
		return 0, err
	}
	if exist {
		log.Count++
		return log.Count, srv.scanRepo.Save(log)
	}
	log = &entity.ScanLog{
		ImageUrl:   imageUrl,
		Hash:       imageHash,
		Count:      1,
		ScanResult: scanResult,
	}
	return log.Count, srv.scanRepo.Create(log)
}
func (srv OCRService) GetImageScanCount(imageHash string) (int, error) {
	log, exist, err := srv.scanRepo.FindByHash(imageHash)
	if err != nil {
		return 0, err
	}
	if exist {
		return log.Count, nil
	}
	return 0, nil
}

// CheckImageScanCount 检查同一个hash图片扫描次数是否达到上限 并且返回图片的hash值
func (srv OCRService) CheckImageScanCount(imageUrl string, maxCount int) (hash string, err error) {
	hash, err = srv.GetImageHash(imageUrl)
	if err != nil {
		return "", err
	}

	count, err := srv.GetImageScanCount(hash)
	if err != nil {
		return "", err
	}

	if count >= maxCount {
		return "", errno.ErrCommon.WithMessage("重复扫描")
	}

	return hash, nil
}
