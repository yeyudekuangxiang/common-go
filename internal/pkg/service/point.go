package service

import (
	"errors"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	messageSrv "mio/internal/pkg/service/message"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/validator"
	"mio/pkg/errno"
	"time"
)

func NewPointService(ctx *context.MioContext) PointService {
	return PointService{ctx: ctx, repo: repository.NewPointRepository(ctx), repoInvite: repository.NewInviteRepository(ctx)}
}

type PointService struct {
	ctx        *context.MioContext
	repo       *repository.PointRepository
	repoInvite *repository.InviteRepository
}

// FindByUserId 获取用户积分
func (srv PointService) FindByUserId(userId int64) (*entity.Point, error) {
	user, err := DefaultUserService.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user.OpenId == "" {
		return &entity.Point{}, errno.ErrUserNotFound
	}
	return srv.FindByOpenId(user.OpenId)
}

// FindByOpenId 获取用户积分
func (srv PointService) FindByOpenId(openId string) (*entity.Point, error) {
	if openId == "" {
		return &entity.Point{}, errno.ErrUserNotFound
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
		go srv.trackPoint(dto, "操作频繁")
		return 0, errors.New("操作频繁")
	}
	defer util.DefaultLock.UnLock(lockKey)

	//检测积分发放次数限制
	if dto.ChangePoint >= 0 {
		limitService := NewPointTransactionCountLimitService(srv.ctx)
		err := limitService.CheckLimitAndUpdate(dto.Type, dto.OpenId)
		if err != nil {
			//积分打点
			go srv.trackPoint(dto, err.Error())
			return 0, err
		}
	}

	//发放或扣减积分
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
			return errors.New("积分不足")
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

		//积分变动需要提醒的type

		//var types = []entity.PointTransactionType{entity.POINT_STEP, entity.POINT_BIKE_RIDE, entity.POINT_COFFEE_CUP, entity.POINT_ECAR}

		//发小程序订阅消息
		message := messageSrv.MiniChangePointTemplate{
			Point:    dto.ChangePoint,
			Source:   dto.Type.Text(),
			Time:     time.Now().Format("2006年01月02日"),
			AllPoint: point.Balance,
		}
		service := messageSrv.MessageService{}
		_, messageErr := service.SendMiniSubMessage(dto.OpenId, config.MessageJumpUrls.ChangePoint, message)
		if messageErr != nil {

		}

		//发完积分，更新邀请表发积分状态
		if dto.InviteId != 0 && dto.Type == entity.POINT_INVITE {
			//更新状态
			err = srv.repoInvite.UpdateIsReward(dto.InviteId)
			if err != nil {
				return err
			}
		}
		return nil
	})

	//积分打点
	if err != nil {
		go srv.trackPoint(dto, err.Error())
	} else {
		go srv.trackPoint(dto, "")
	}

	return balance, err
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
			return errors.New("积分不足")
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
		return errno.ErrUserNotFound
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

func (srv PointService) trackPoint(dto srv_types.ChangeUserPointDTO, failMessage string) {
	DefaultZhuGeService().TrackPoint(srv_types.TrackPoint{
		OpenId:      dto.OpenId,
		PointType:   dto.Type,
		ChangeType:  util.Ternary(dto.ChangePoint > 0, "inc", "dec").String(),
		Value:       uint(dto.ChangePoint),
		IsFail:      util.Ternary(failMessage == "", false, true).Bool(),
		FailMessage: failMessage,
	})
}
