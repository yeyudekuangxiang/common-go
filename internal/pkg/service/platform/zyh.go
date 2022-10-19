package platform

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/config"
	"mio/internal/pkg/core/context"
	entityV2 "mio/internal/pkg/model/entity"
	entity "mio/internal/pkg/model/entity/activity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/activity"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
	util2 "mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	"mio/pkg/errno"
)

func NewZyhService(ctx *context.MioContext) *ZyhService {
	return &ZyhService{
		ctx:              ctx,
		Domain:           config.Config.ActivityZyh.Domain,
		ZyhRepository:    activity.NewZyhRepository(ctx),
		ZyhLogRepository: activity.NewZyhLogRepository(ctx),
		UserRepository:   repository.NewUserRepository(),
	}
}

type ZyhService struct {
	ctx              *context.MioContext
	Domain           string `json:"domain"`
	ZyhRepository    activity.ZyhRepository
	ZyhLogRepository activity.ZyhLogRepository
	UserRepository   repository.UserRepository
}

func (srv ZyhService) CheckIsVolunteer(openid string) (bool, error) {
	info := srv.ZyhRepository.FindBy(activity.FindZyhById{
		Openid: openid,
	})
	if info.Id == 0 {
		return false, errno.ErrCommon.WithMessage("不存在该志愿者")
	}
	return true, nil
}

//志愿汇发金币

func (srv ZyhService) SendPoint(pointType string, openid string, point string) (code string, error error) {
	info := srv.ZyhRepository.FindBy(activity.FindZyhById{
		Openid: openid,
	})
	if info.Id == 0 {
		return "30000", errno.ErrCommon.WithMessage("不存在该志愿者")
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
		return response.ErrCode, errno.ErrCommon.WithMessage(response.Message)
	}
	//入库
	fmt.Printf("%v\n", response)
	return response.ErrCode, nil
}

func (srv ZyhService) SendPointJb(pointType string, volId string, point string) (code string, error error) {
	params := make(map[string]string, 0)
	params["AccessKeyId"] = config.Config.ActivityZyh.AccessKeyId
	params["type"] = pointType
	params["volunteerId"] = volId
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
		return response.ErrCode, errno.ErrCommon.WithMessage(response.Message)
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
	return errno.ErrCommon.WithMessage("志愿汇用户已存在")
}

func (srv ZyhService) GetZyhInfoByMobile(dto srv_types.GetZyhOpenDTO) (gin.H, error) {
	openid := ""
	nick := ""
	if dto.Mobile != "" {
		info, exit, err := srv.UserRepository.GetUser(repository.GetUserBy{
			Mobile: dto.Mobile,
		})
		if !exit {
			return gin.H{
				"结果": "该手机号未注册小程序",
			}, nil
		}
		if err != nil {
			return gin.H{
				"结果": err.Error(),
			}, err
		}
		openid = info.OpenId
		nick = info.Nickname
		if openid == "" {
			return gin.H{
				"结果": "已注册绿喵小程序，但openid为空",
			}, err
		}
		ret, err := srv.GetInfoBy(srv_types.GetZyhGetInfoByDTO{Openid: openid})
		if err != nil {
			return gin.H{
				"结果": err.Error(),
				"昵称": nick,
			}, nil
		}
		if ret.Id == 0 {
			return gin.H{
				"手机号": info.PhoneNumber,
				"结果":  "手机号未绑定志愿者id",
				"昵称":  nick,
			}, nil
		}
		list, err := srv.ZyhLogRepository.GetListBy(repotypes.GetZyhListBy{Openid: openid})
		if err != nil {
			return gin.H{
				"手机号":   info.PhoneNumber,
				"昵称":    nick,
				"志愿者编号": ret.VolId,
				"结果":    err.Error(),
			}, nil
		}
		createList := make([]srv_types.GetZyhLogDTO, 0)
		for _, l := range list {
			sendType := "其他"
			switch l.PointType {
			case entityV2.POINT_QUIZ:
				sendType = "答题"
				break
			case entityV2.POINT_STEP:
				sendType = "步行"
				break
			}
			createList = append(createList, srv_types.GetZyhLogDTO{
				PointType:  sendType,
				PointValue: l.Value,
				ResultCode: l.ResultCode,
				CreateTime: l.CreatedAt.Format("2006.01.02 15:04:05"),
			})
		}
		return gin.H{
			"昵称":    nick,
			"手机号":   info.PhoneNumber,
			"志愿者编号": ret.VolId,
			"积分明细":  createList,
		}, nil
	}
	return gin.H{
		"结果": "手机号为空",
	}, nil
}

func (srv ZyhService) GetZyhInfoByVolId(dto srv_types.GetZyhOpenDTO) (gin.H, error) {
	if dto.VolId != "" {
		openid := ""
		nick := ""
		info, err := srv.GetInfoBy(srv_types.GetZyhGetInfoByDTO{VolId: dto.VolId})
		if err != nil {
			return gin.H{
				"结果": err.Error(),
				"昵称": nick,
			}, nil
		}
		if info.Id == 0 {
			return gin.H{
				"结果": "志愿者id未绑定关系",
				"昵称": nick,
			}, nil
		}

		if info.Openid != "" {
			userInfo, exit, UserErr := srv.UserRepository.GetUser(repository.GetUserBy{
				OpenId: info.Openid,
			})
			if !exit {
				return gin.H{
					"结果": "该手机号未注册小程序",
				}, nil
			}
			if UserErr != nil {
				return gin.H{
					"结果": UserErr.Error(),
				}, UserErr
			}
			openid = userInfo.OpenId
			nick = userInfo.Nickname

			list, err := srv.ZyhLogRepository.GetListBy(repotypes.GetZyhListBy{Openid: openid})
			if err != nil {
				return gin.H{
					"结果": err.Error(),
					"昵称": nick,
				}, nil
			}
			createList := make([]srv_types.GetZyhLogDTO, 0)
			for _, l := range list {
				sendType := "其他"
				switch l.PointType {
				case entityV2.POINT_QUIZ:
					sendType = "答题"
					break
				case entityV2.POINT_STEP:
					sendType = "步行"
					break
				}
				createList = append(createList, srv_types.GetZyhLogDTO{
					PointType:  sendType,
					PointValue: l.Value,
					ResultCode: l.ResultCode,
					CreateTime: l.CreatedAt.Format("2006.01.02 15:04:05"),
				})
			}
			return gin.H{
				"昵称":    nick,
				"手机号":   userInfo.PhoneNumber,
				"志愿者编号": info.VolId,
				"积分明细":  createList,
			}, nil
		}
	}

	return gin.H{}, nil

}
