package platform

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"mio/config"
	"mio/internal/pkg/core/context"
	entity "mio/internal/pkg/model/entity/activity"
	"mio/internal/pkg/repository/activity"
	"mio/internal/pkg/service/srv_types"
	util2 "mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
)

func NewZyhService(ctx *context.MioContext) *ZyhService {
	return &ZyhService{
		ctx:              ctx,
		Domain:           config.Config.ActivityZyh.Domain,
		ZyhRepository:    activity.NewZyhRepository(ctx),
		ZyhLogRepository: activity.NewZyhLogRepository(ctx),
	}
}

type ZyhService struct {
	ctx              *context.MioContext
	Domain           string `json:"domain"`
	ZyhRepository    activity.ZyhRepository
	ZyhLogRepository activity.ZyhLogRepository
}

//志愿汇发金币

func (srv ZyhService) SendPoint(pointType string, openid string, point string) (code string, error error) {
	info := srv.ZyhRepository.FindBy(activity.FindZyhById{
		Openid: openid,
	})
	if info.Id == 0 {
		return "30000", errors.New("不存在该志愿者")
	}
	params := make(map[string]string, 0)
	params["AccessKeyId"] = config.Config.ActivityZyh.AccessKeyId
	params["type"] = pointType
	params["volunteerId"] = info.VolId
	params["point"] = point
	sign := util2.GetSign(params)
	params["Signature"] = sign //strings.ToUpper()
	url := srv.Domain
	body, err := httputil.PostMapFrom(url, params)
	fmt.Printf("%s\n", body)
	if err != nil {
		return "30001", err
	}
	response := zyhCommonResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "30002", err
	}
	if response.ErrCode != "0000" {
		return response.ErrCode, errors.New(response.Message)
	}
	//入库
	fmt.Printf("%v\n", response)
	return response.ErrCode, nil
}

func (srv ZyhService) CreateLog(dto srv_types.GetZyhLogAddDTO) error {
	return srv.ZyhLogRepository.Save(&entity.ZyhLog{
		Openid:         dto.Openid,
		PointType:      dto.PointType,
		Value:          dto.Value,
		ResultCode:     dto.ResultCode,
		AdditionalInfo: dto.AdditionalInfo,
		TransactionId:  dto.TransactionId,
	})
}

func (srv ZyhService) GetInfoBy(dto srv_types.GetZyhGetInfoByDTO) (entity.Zyh, error) {
	info := srv.ZyhRepository.FindBy(activity.FindZyhById{
		Openid: dto.Openid,
		VolId:  dto.VolId,
	})
	return info, nil
}

func (srv ZyhService) Create(dto srv_types.GetZyhGetInfoByDTO) error {
	info, _ := srv.GetInfoBy(dto)
	if info.Id == 0 {
		//入库
		return srv.ZyhRepository.Save(&entity.Zyh{
			Openid: dto.Openid,
			VolId:  dto.VolId,
		})
	}
	return errors.New("志愿汇用户已存在")
}
