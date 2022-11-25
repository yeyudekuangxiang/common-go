package service

import (
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/point/cmd/rpc/pointclient"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	messageSrv "mio/internal/pkg/service/message"
	platformSrv "mio/internal/pkg/service/platform"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/validator"
	"mio/pkg/errno"
	"strconv"
	"time"
)

func NewPointService(ctx *context.MioContext) PointService {
	return PointService{ctx: ctx, repo: repository.NewPointRepository(ctx), repoInvite: repository.NewInviteRepository(ctx), repoPointTransaction: repository.NewPointTransactionRepository(ctx)}
}

type PointService struct {
	ctx                  *context.MioContext
	repo                 *repository.PointRepository
	repoInvite           *repository.InviteRepository
	repoPointTransaction *repository.PointTransactionRepository
}

// FindByUserId 获取用户积分
func (srv PointService) FindByUserId(userId int64) (*entity.Point, error) {
	user, err := DefaultUserService.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user.OpenId == "" {
		return &entity.Point{}, errno.ErrUserNotFound.WithCaller()
	}
	return srv.FindByOpenId(user.OpenId)
}

// FindByOpenId 获取用户积分
func (srv PointService) FindByOpenId(openId string) (*entity.Point, error) {
	if openId == "" {
		return &entity.Point{}, errno.ErrUserNotFound.WithCaller()
	}
	point := srv.repo.FindBy(repository.FindPointBy{
		OpenId: openId,
	})
	return &point, nil
}

// IncUserPoint 加积分操作
func (srv PointService) IncUserPoint(dto srv_types.IncUserPointDTO) (int64, error) {
	changePointDto := srv_types.ChangeUserPointDTO{}
	if err := util.MapTo(dto, &changePointDto); err != nil {
		return 0, err
	}
	return srv.changeUserPoint(changePointDto)
}

// DecUserPoint 减积分操作
func (srv PointService) DecUserPoint(dto srv_types.DecUserPointDTO) (int64, error) {
	if dto.ChangePoint < 0 {
		return 0, errors.New("DecUserPoint Value error")
	}
	changePointDto := srv_types.ChangeUserPointDTO{}
	if err := util.MapTo(dto, &changePointDto); err != nil {
		return 0, err
	}
	changePointDto.ChangePoint = -changePointDto.ChangePoint
	return srv.changeUserPoint(changePointDto)
}

//changeUserPoint 变动积分操作
func (srv PointService) changeUserPoint(dto srv_types.ChangeUserPointDTO) (int64, error) {

	lockKey := "changeUserPoint" + dto.OpenId
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return 0, errno.ErrCommon.WithMessage("操作频繁")
	}
	defer util.DefaultLock.UnLock(lockKey)

	//检测积分发放次数限制
	if dto.ChangePoint >= 0 {
		limitService := NewPointTransactionCountLimitService(srv.ctx)
		err := limitService.CheckLimitAndUpdate(dto.Type, dto.OpenId)
		if err != nil {
			return 0, err
		}
	}

	var balance int64
	var err error
	var ptId int64
	if dto.ChangePoint > 0 {
		resp, err := app.RpcService.PointRpcSrv.IncPoint(srv.ctx, &pointclient.IncPointReq{
			Openid:       dto.OpenId,
			Type:         string(dto.Type),
			BizId:        dto.BizId,
			ChangePoint:  uint64(dto.ChangePoint),
			AdminId:      uint64(dto.AdminId),
			Node:         dto.Note,
			AdditionInfo: dto.AdditionInfo,
			InviteId:     dto.InviteId,
		})
		if err != nil {
			return 0, err
		}
		balance = resp.Point
		ptId = resp.PointTransactionId
	} else {
		resp, err := app.RpcService.PointRpcSrv.DecPoint(srv.ctx, &pointclient.DecPointReq{
			Openid:       dto.OpenId,
			Type:         string(dto.Type),
			BizId:        dto.BizId,
			ChangePoint:  uint64(-dto.ChangePoint),
			AdminId:      uint64(dto.AdminId),
			Node:         dto.Note,
			AdditionInfo: dto.AdditionInfo,
			InviteId:     dto.InviteId,
		})
		if err != nil {
			return 0, err
		}
		balance = resp.Point
		ptId = resp.PointTransactionId
	}
	go srv.afterChangePoint(ptId, balance, dto)
	return balance, err
}
func (srv PointService) afterChangePoint(ptId int64, balance int64, dto srv_types.ChangeUserPointDTO) {
	//积分变动提醒
	types := map[entity.PointTransactionType]string{
		entity.POINT_JHX:                    "金华行",
		entity.POINT_FAST_ELECTRICITY:       "快电",
		entity.POINT_ECAR:                   "星星充电",
		entity.POINT_RECYCLING_CLOTHING:     "旧物回收 oola衣物鞋帽",
		entity.POINT_RECYCLING_DIGITAL:      "旧物回收 oola数码",
		entity.POINT_RECYCLING_APPLIANCE:    "旧物回收 oola家电",
		entity.POINT_RECYCLING_BOOK:         "旧物回收 oola书籍",
		entity.POINT_FMY_RECYCLING_CLOTHING: "旧物回收 fmy衣物鞋帽",
		entity.POINT_SYSTEM_ADD:             "系统补发",
	}
	_, ok := types[dto.Type]
	if ok && (dto.ChangePoint > 0) {
		//发小程序订阅消息
		message := messageSrv.MiniChangePointTemplate{
			Point:    dto.ChangePoint,
			Source:   dto.Type.Text(),
			Time:     time.Now().Format("2006年01月02日"),
			AllPoint: balance,
		}
		service := messageSrv.MessageService{}
		_, messageErr := service.SendMiniSubMessage(dto.OpenId, config.MessageJumpUrls.ChangePoint, message)
		if messageErr != nil {

		}
	}

	//同步到志愿汇
	if dto.ChangePoint >= 0 {
		//积分变动提醒
		typeZyh := map[entity.PointTransactionType]string{
			entity.POINT_STEP:                    "步行",
			entity.POINT_COFFEE_CUP:              "自带咖啡杯",
			entity.POINT_BIKE_RIDE:               "骑行",
			entity.POINT_ECAR:                    "答题活动",
			entity.POINT_QUIZ:                    "答题活动",
			entity.POINT_JHX:                     "金华行",
			entity.POINT_POWER_REPLACE:           "电车换电",
			entity.POINT_DUIBA_INTEGRAL_RECHARGE: "兑吧虚拟商品充值积分",
			entity.POINT_RECYCLING_CLOTHING:      "旧物回收 oola衣物鞋帽",
			entity.POINT_RECYCLING_DIGITAL:       "旧物回收 oola数码",
			entity.POINT_RECYCLING_APPLIANCE:     "旧物回收 oola家电",
			entity.POINT_RECYCLING_BOOK:          "旧物回收 oola书籍",
			entity.POINT_FMY_RECYCLING_CLOTHING:  "旧物回收 fmy衣物鞋帽",
			entity.POINT_FAST_ELECTRICITY:        "快电",
			entity.POINT_REDUCE_PLASTIC:          "环保减塑",
		}
		_, zyhOk := typeZyh[dto.Type]
		if dto.Type == entity.POINT_QUIZ && dto.ChangePoint == 2500 {
			zyhOk = false
		}
		if zyhOk {
			sendType := "0"
			switch dto.Type {
			case entity.POINT_QUIZ:
				sendType = "1"
				break
			case entity.POINT_STEP:
				sendType = "2"
				break
			}
			serviceZyh := platformSrv.NewZyhService(context.NewMioContext())
			messageCode, messageErr := serviceZyh.SendPoint(sendType, dto.OpenId, strconv.FormatInt(dto.ChangePoint, 10))
			if messageCode != "30000" {
				//发送结果记录到日志
				msgErr := ""
				if messageErr != nil {
					msgErr = messageErr.Error()
				}
				serviceZyh.CreateLog(srv_types.GetZyhLogAddDTO{
					Openid:         dto.OpenId,
					PointType:      dto.Type,
					Value:          dto.ChangePoint,
					ResultCode:     messageCode,
					AdditionalInfo: msgErr,
					TransactionId:  dto.BizId,
				})
			}
		}
	}

	//发完积分，更新邀请表发积分状态
	if dto.InviteId != 0 && dto.Type == entity.POINT_INVITE {
		//更新状态
		err := srv.repoInvite.UpdateIsReward(dto.InviteId)
		if err != nil {
			app.Logger.Error("更新邀请状态失败", dto.InviteId, err)
		}
	}
}

//ChangeUserPointByOffline 线下发积分
func (srv PointService) ChangeUserPointByOffline(dto srv_types.ChangeUserPointDTO) (int64, error) {
	var balance int64 = 0
	err := srv.ctx.Transaction(func(ctx *context.MioContext) error {
		//查询积分账户
		pointRepo := repository.NewPointRepository(ctx)
		point, err := pointRepo.FindForUpdate(dto.OpenId)
		if err != nil {
			return err
		}
		//判读积分余额是否充足
		if dto.ChangePoint < 0 && point.Balance+dto.ChangePoint < 0 {
			return errno.ErrCommon.WithMessage("积分不足")
		}
		if point.Id == 0 {
			//创建积分账户
			point.OpenId = dto.OpenId
			point.Balance += dto.ChangePoint
			if err := pointRepo.Create(&point); err != nil {
				return err
			}
		} else {
			//更新积分账户
			point.Balance += dto.ChangePoint
			if err := pointRepo.Save(&point); err != nil {
				return err
			}
		}
		balance = point.Balance
		//增加积分变动记录
		tranService := NewPointTransactionService(ctx)
		_, err = tranService.CreateTransaction(CreatePointTransactionParam{
			BizId:        dto.BizId,
			OpenId:       dto.OpenId,
			Type:         dto.Type,
			Value:        dto.ChangePoint,
			AdminId:      dto.AdminId,
			Note:         dto.Note,
			AdditionInfo: dto.AdditionInfo,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return balance, err
}

//AdminAdjustUserPoint 管理员变动积分
func (srv PointService) AdminAdjustUserPoint(adminId int, param AdminAdjustUserPointParam) error {

	if err := validator.ValidatorStruct(param); err != nil {
		app.Logger.Error(param, err)
		return err
	}

	user, err := DefaultUserService.GetUserBy(repository.GetUserBy{
		OpenId:     param.OpenId,
		LikeMobile: param.Phone,
	})
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errno.ErrUserNotFound.WithCaller()
	}
	value := param.Value
	if param.Type == entity.POINT_SYSTEM_REDUCE {
		value = -value
	}

	_, err = srv.changeUserPoint(srv_types.ChangeUserPointDTO{
		OpenId:      param.OpenId,
		Type:        param.Type,
		ChangePoint: value,
		AdminId:     adminId,
		Note:        param.Note,
		BizId:       util.UUID(),
	})
	return err
}
