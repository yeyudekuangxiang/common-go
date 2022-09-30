package platform

import (
	"encoding/json"
	"errors"
	"fmt"
	"mio/config"
	"mio/internal/pkg/core/context"
	entity "mio/internal/pkg/model/entity/activity"
	"mio/internal/pkg/repository/activity"
	util2 "mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
)

func NewZyhService(ctx *context.MioContext) *ZyhService {
	return &ZyhService{
		ctx:              ctx,
		Domain:           "http://47.99.112.147:8080",
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

func (srv ZyhService) SendPoint(pointType string, openid string, point string) (code string, error error) {
	info := srv.ZyhRepository.FindBy(activity.FindZyhById{
		Openid: openid,
	})
	if info.Id == 0 {
		return "30000", errors.New("不存在该志愿者")
	}
	volunteerId := info.Openid
	params := make(map[string]string, 0)
	params["AccessKeyId"] = config.Config.ActivityZyh.AccessKeyId
	params["type"] = pointType
	params["volunteerId"] = volunteerId
	params["point"] = point
	sign := util2.GetSign(params)
	params["Signature"] = sign //strings.ToUpper()
	url := srv.Domain + "/VMSAPI/api/publicSchool/lm/pushPointsInfo.do"
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

func (srv ZyhService) SendPointLog(log entity.ZyhLog) error {
	ret := srv.ZyhLogRepository.Save(&log)
	return ret
}
