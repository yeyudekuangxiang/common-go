package service

import (
	"errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"time"
)

func NewPointService(ctx *context.MioContext) PointService {
	return PointService{ctx: ctx, repo: repository.NewPointRepository(ctx)}
}

type PointService struct {
	ctx  *context.MioContext
	repo repository.PointRepository
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
		return 0, errors.New("操作频繁")
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
		return nil
	})
	return balance, err
}

//AdminAdjustUserPoint 管理员变动积分
func (srv PointService) AdminAdjustUserPoint(adminId int, param AdminAdjustUserPointParam) error {

	if err := util.ValidatorStruct(param); err != nil {
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
	})
	return err
}
