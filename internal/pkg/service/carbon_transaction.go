package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"mio/config"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"strconv"
	"time"
)

var DefaultCarbonTransactionService = NewCarbonTransactionService(repository.DefaultCarbonTransactionRepository)

func NewCarbonTransactionService(repo repository.CarbonTransactionRepository) CarbonTransactionService {
	return CarbonTransactionService{
		repo: repo,
	}
}

type CarbonTransactionService struct {
	repo repository.CarbonTransactionRepository
}

// Create 添加发放碳量记录并且更新用户剩余碳量
func (srv CarbonTransactionService) Create(dto api_types.CreateCarbonTransactionDto) (*entity.CarbonTransaction, error) {
	if err := util.ValidatorStruct(dto); err != nil {
		app.Logger.Error(dto, err)
		return nil, err
	}
	/*err := DefaultPointTransactionCountLimitService.CheckLimitAndUpdate(param.Type, param.OpenId)
	if err != nil {
		return nil, err
	}*/
	//入库
	bannerDo := entity.CarbonTransaction{
		TransactionId: util.UUID(),
		Info:          dto.Info,
		CreatedAt:     model.Time{Time: time.Now()},
		UpdatedAt:     model.Time{Time: time.Now()}}
	if err := util.MapTo(dto, &bannerDo); err != nil {
		return nil, err
	}
	err := srv.repo.Create(&bannerDo)
	if err != nil {
		return nil, err
	}
	//记录redis,今日榜单
	todayDate := time.Now().Format("20060102")        //当天时间 年月日
	score := dto.Value                                //增加的碳量
	UserIdString := strconv.FormatInt(dto.UserId, 10) //用户uid
	redisKey := fmt.Sprintf(config.RedisKey.UserCarbonRank, todayDate)
	err1 := app.Redis.ZIncrBy(context.Background(), redisKey, score, UserIdString).Err()
	if err1 != nil {
		return nil, err1
	}
	//更新用户表，碳量字段
	DefaultUserService.UpdateUserCarbon(dto.UserId, score)
	return nil, nil
}

// Bank 排行榜
func (srv CarbonTransactionService) Bank(dto api_types.GetCarbonTransactionBankDto) ([]api_types.CarbonTransactionBank, error) {
	//1 获取排行榜信息
	todayDate := time.Now().Format("20060102") //当天时间 年月日
	redisKey := fmt.Sprintf(config.RedisKey.UserCarbonRank, todayDate)

	op := &redis.ZRangeBy{
		Offset: dto.Offset, // 类似sql的limit, 表示开始偏移量
		Count:  dto.Limit,  // 一次返回多少数据
	}
	bank, err := app.Redis.ZRangeByScoreWithScores(context.Background(), redisKey, op).Result()
	if err != nil {
		return nil, err
	}

	//2 排行榜当前页所有的uid
	uids := make([]int64, len(bank))
	for _, val := range bank {
		uids = append(uids, val.Member.(int64))
	}

	//3 根据uid获取用户信息
	userList, _ := DefaultUserService.GetUserListBy(repository.GetUserListBy{UserIds: uids})
	var userMap map[int64]entity.User /*创建用户集合 */
	for _, val := range userList {
		userMap[val.ID] = entity.User{
			Nickname:  val.Nickname,
			AvatarUrl: val.AvatarUrl,
		}
	}
	//4 todo 根据当前页的uid和我的uid去好友表查到存在的uid

	//5 根据用户信息和好友信息进行信息整理
	CarbonBankList := make([]api_types.CarbonTransactionBank, 0)
	for key, val := range bank {
		CarbonBankList = append(CarbonBankList, api_types.CarbonTransactionBank{
			Carbon: val.Score,      //分数
			Rank:   int64(key) + 1, //排行
			User: api_types.CarbonUser{
				AvatarUrl: userMap[val.Member.(int64)].AvatarUrl,
				Nickname:  userMap[val.Member.(int64)].Nickname,
				Uid:       val.Member.(int64), //用户uid
			},
			IsFriend: false, // todo 是否是好友
		})
	}
	return CarbonBankList, nil
}

// MyBank 我的排名
func (srv CarbonTransactionService) MyBank(dto api_types.GetCarbonTransactionMyBankDto) (api_types.CarbonTransactionMyBank, error) {
	todayDate := time.Now().Format("20060102") //当天时间 年月日
	redisKey := fmt.Sprintf(config.RedisKey.UserCarbonRank, todayDate)

	UserIdString := strconv.FormatInt(dto.UserId, 10) //我的uid string

	mySort := app.Redis.ZRank(context.Background(), redisKey, UserIdString)    //我的当前排名
	myCarbon := app.Redis.ZScore(context.Background(), redisKey, UserIdString) //我的碳量
	allCount := app.Redis.ZCard(context.Background(), redisKey)                //总人数

	mySortDec := decimal.NewFromInt(mySort.Val())
	allCountDec := decimal.NewFromInt(allCount.Val())

	//超越**%用户
	var overPer string
	overPer = mySortDec.Div(allCountDec).Round(2).Mul(decimal.NewFromInt(100)).String() + "%"

	//我的当前排名
	myRank := allCountDec.Sub(allCountDec)

	//我的信息
	user, _ := DefaultUserService.GetUserById(dto.UserId)

	myBank := api_types.CarbonTransactionMyBank{
		OverPer: overPer,
		Carbon:  myCarbon.Val(),
		Rank:    myRank.String(),
		User: api_types.CarbonUser{
			AvatarUrl: user.AvatarUrl, //头像
			Nickname:  user.Nickname}, //昵称
	}
	return myBank, nil
}

// Info 我的减碳成就-基础信息
func (srv CarbonTransactionService) Info(dto api_types.GetCarbonTransactionInfoDto) {

}

// Classify 我的减碳成就-我的减碳足迹
func (srv CarbonTransactionService) Classify(dto api_types.GetCarbonTransactionClassifyDto) {

}

// History 我的减碳成就-近2周减碳
func (srv CarbonTransactionService) History(dto api_types.GetCarbonTransactionHistoryDto) {

}
