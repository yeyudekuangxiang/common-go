package service

import (
	contextRedis "context"
	"fmt"
	"github.com/pkg/errors"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	carbonProducer "mio/internal/pkg/queue/producer/carbon"
	carbonmsg "mio/internal/pkg/queue/types/message/carbon"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"strconv"
	"time"
)

func NewCarbonService(ctx *context.MioContext) CarbonService {
	return CarbonService{ctx: ctx,
		repo:      repository.NewCarbonRepository(ctx),
		repoScene: repository.NewCarbonSceneRepository(ctx),
		repoT:     repository.NewCarbonTransactionRepository(ctx)}
}

type CarbonService struct {
	ctx       *context.MioContext
	repo      repository.CarbonRepository
	repoScene repository.CarbonSceneRepository
	repoT     repository.CarbonTransactionRepository
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
		return 0, errno.ErrCommon.WithMessage("操作频繁")
	}
	defer util.DefaultLock.UnLock(lockKey)
	//发放或扣减积分
	var balance float64 = 0
	err := srv.ctx.Transaction(func(ctx *context.MioContext) error {
		//查询积分账户
		carbonRepo := repository.NewCarbonRepository(ctx)
		carbon, err := carbonRepo.FindForUpdate(dto.OpenId)
		if err != nil {
			return err
		}

		//判读积分余额是否充足
		if dto.ChangePoint < 0 && carbon.Carbon+dto.ChangePoint < 0 {
			return errno.ErrCommon.WithMessage("碳量不足")
		}

		if carbon.Id == 0 {
			//创建积分账户
			carbon.OpenId = dto.OpenId
			carbon.Carbon += dto.ChangePoint
			if err := carbonRepo.Create(&carbon); err != nil {
				return err
			}
		} else {
			//更新积分账户
			carbon.Carbon += dto.ChangePoint
			if err := carbonRepo.Save(&carbon); err != nil {
				return err
			}
		}

		//入库记录表
		CarbonTransactionDo := entity.CarbonTransaction{
			OpenId:        dto.OpenId,
			UserId:        dto.Uid,
			Type:          dto.Type,
			Info:          dto.AdditionInfo,
			City:          dto.CityCode,
			Value:         dto.ChangePoint,
			AdminId:       dto.AdminId,
			TransactionId: util.UUID(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		err = srv.repoT.Create(&CarbonTransactionDo)
		if err != nil {
			return err
		}
		//投递mq
		if err := carbonProducer.ChangeSuccessToQueue(carbonmsg.CarbonChangeSuccess{
			Openid:        dto.OpenId,
			UserId:        dto.Uid,
			TransactionId: CarbonTransactionDo.TransactionId,
			Type:          string(dto.Type),
			City:          dto.CityCode,
			Value:         dto.ChangePoint,
			Info:          dto.AdditionInfo,
		}); err != nil {
			app.Logger.Errorf("ChangeSuccessToQueue 投递失败:%v", err)
		}
		//记录redis,今日榜单
		UserIdString := strconv.FormatInt(dto.Uid, 10) //用户uid
		redisKey := fmt.Sprintf(config.RedisKey.UserCarbonRank, time.Now().Format("20060102"))
		errRedis := app.Redis.ZIncrBy(contextRedis.Background(), redisKey, dto.ChangePoint, UserIdString).Err()
		if errRedis != nil {
			return errRedis
		}
		balance = carbon.Carbon
		return nil
	})
	return balance, err
}
