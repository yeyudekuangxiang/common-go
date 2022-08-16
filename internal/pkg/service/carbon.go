package service

import (
	"errors"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"time"
)

func NewCarbonService(ctx *context.MioContext) CarbonService {
	return CarbonService{ctx: ctx, repo: repository.NewCarbonRepository(ctx)}
}

type CarbonService struct {
	ctx  *context.MioContext
	repo repository.CarbonRepository
}

// FindByUserId 获取用户碳量
func (srv CarbonService) FindByUserId(userId int64) (*entity.Carbon, error) {
	user, err := DefaultUserService.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user.OpenId == "" {
		return &entity.Carbon{}, errno.ErrUserNotFound
	}
	return srv.FindByOpenId(user.OpenId)
}

// FindByOpenId 获取用户碳量
func (srv CarbonService) FindByOpenId(openId string) (*entity.Carbon, error) {
	if openId == "" {
		return &entity.Carbon{}, errno.ErrUserNotFound
	}
	point := srv.repo.FindBy(repository.FindCarbonBy{
		OpenId: openId,
	})
	return &point, nil
}

// IncUserCarbon 加碳量操作
func (srv CarbonService) IncUserCarbon(dto srv_types.IncUserCarbonDTO) (float64, error) {
	changePointDto := srv_types.ChangeUserCarbonDTO{}
	if err := util.MapTo(dto, &changePointDto); err != nil {
		return 0, err
	}
	return srv.changeUserPoint(changePointDto)
}

// DecUserCarbon 减积分操作
func (srv CarbonService) DecUserCarbon(dto srv_types.DecUserPointDTO) (float64, error) {
	if dto.ChangePoint < 0 {
		return 0, errors.New("DecUserPoint Value error")
	}
	changePointDto := srv_types.ChangeUserCarbonDTO{}
	if err := util.MapTo(dto, &changePointDto); err != nil {
		return 0, err
	}
	changePointDto.ChangePoint = -changePointDto.ChangePoint
	return srv.changeUserPoint(changePointDto)
}

//changeUserPoint 变动积分操作
func (srv CarbonService) changeUserPoint(dto srv_types.ChangeUserCarbonDTO) (float64, error) {
	lockKey := "changeUserCarbon" + dto.OpenId
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return 0, errors.New("操作频繁")
	}
	defer util.DefaultLock.UnLock(lockKey)

	//检测积分发放次数限制
	/*if dto.ChangePoint >= 0 {
		limitService := NewPointTransactionCountLimitService(srv.ctx)
		err := limitService.CheckLimitAndUpdate(dto.Type, dto.OpenId)
		if err != nil {
			return 0, err
		}
	}*/

	//发放或扣减积分
	var balance float64 = 0
	err := srv.ctx.Transaction(func(ctx *context.MioContext) error {

		//查询积分账户
		pointRepo := repository.NewCarbonRepository(ctx)
		point, err := pointRepo.FindForUpdate(dto.OpenId)
		if err != nil {
			return err
		}

		//判读积分余额是否充足
		if dto.ChangePoint < 0 && point.Carbon+dto.ChangePoint < 0 {
			return errors.New("积分不足")
		}

		if point.Id == 0 {
			//创建积分账户
			point.OpenId = dto.OpenId
			point.Carbon += dto.ChangePoint
			if err := pointRepo.Create(&point); err != nil {
				return err
			}
		} else {
			//更新积分账户
			point.Carbon += dto.ChangePoint
			if err := pointRepo.Save(&point); err != nil {
				return err
			}
		}
		balance = point.Carbon
		return nil
	})
	return balance, err
}
