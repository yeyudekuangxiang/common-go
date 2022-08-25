package service

import (
	contextRedis "context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"mio/config"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/timeutils"
	"mio/pkg/baidu"
	"sort"
	"strconv"
	"time"
)

func NewCarbonTransactionService(ctx *context.MioContext) CarbonTransactionService {
	return CarbonTransactionService{ctx: ctx,
		repo:      repository.NewCarbonTransactionRepository(ctx),
		repoDay:   repository.NewCarbonTransactionDayRepository(ctx),
		repoScene: repository.NewCarbonSceneRepository(ctx),
		repoPoint: repository.NewPointTransactionRepository(ctx),
	}
}

type CarbonTransactionService struct {
	ctx       *context.MioContext
	repo      repository.CarbonTransactionRepository
	repoDay   repository.CarbonTransactionDayRepository
	repoScene repository.CarbonSceneRepository
	repoPoint *repository.PointTransactionRepository
}

func (srv CarbonTransactionService) PointToCarbon() {

	var types []entity.PointTransactionType
	types = append(types, entity.POINT_STEP, entity.POINT_BIKE_RIDE, entity.POINT_COFFEE_CUP, entity.POINT_ECAR)

	//查询场景配置
	sceneInfoSTEP := srv.repoScene.FindBy(repotypes.CarbonSceneBy{Type: entity.CARBON_STEP})
	sceneInfoBIKERIDE := srv.repoScene.FindBy(repotypes.CarbonSceneBy{Type: entity.CARBON_BIKE_RIDE})
	sceneInfoCOFFEECUP := srv.repoScene.FindBy(repotypes.CarbonSceneBy{Type: entity.CARBON_COFFEE_CUP})
	sceneInfoECAR := srv.repoScene.FindBy(repotypes.CarbonSceneBy{Type: entity.CARBON_ECAR})
	var uid int64 = 70
	Offset := 1
	for true {
		list, _ := srv.repoPoint.GetPageListBy(repository.GetPointTransactionPageListBy{
			Types:  types,
			Offset: Offset,
			Limit:  1000,
		})
		Offset = Offset + 1
		if len(list) == 0 {
			println("跑完了ok")
			return
		}
		continue
		Offset = Offset + 1
		for _, transaction := range list {
			valNew := 0.0
			var a entity.CarbonScene
			var scene entity.CarbonTransactionType
			var info string
			if transaction.Type == entity.POINT_STEP {
				scene = entity.CARBON_STEP
				a = sceneInfoSTEP
				valNew, _ = decimal.NewFromInt(transaction.Value).Mul(decimal.NewFromInt(100)).Float64()
				info = fmt.Sprintf("{time=%v, count=%d}", time.Now(), valNew)
			}
			if transaction.Type == entity.POINT_BIKE_RIDE {
				scene = entity.CARBON_BIKE_RIDE
				a = sceneInfoBIKERIDE
				valNew = 1
				info = string(transaction.AdditionalInfo)
			}
			if transaction.Type == entity.POINT_COFFEE_CUP {
				scene = entity.CARBON_COFFEE_CUP
				a = sceneInfoCOFFEECUP
				valNew = 1
				info = string(transaction.AdditionalInfo)
			}
			if transaction.Type == entity.POINT_ECAR {
				scene = entity.CARBON_ECAR
				a = sceneInfoECAR
				valNew, _ = decimal.NewFromInt(transaction.Value).Div(decimal.NewFromInt(10)).Round(2).Float64()
				info = fmt.Sprintf("{time=%v, count=%d}", time.Now(), valNew)
			}
			carbon := srv.repoScene.GetValue(a, valNew) //增加的碳量
			ret, _ := NewCarbonService(context.NewMioContext()).IncUserCarbon(srv_types.IncUserCarbonDTO{
				OpenId:       transaction.OpenId,
				Type:         scene,
				BizId:        transaction.TransactionId,
				ChangePoint:  carbon,
				AdditionInfo: info,
				CityCode:     "",
				Uid:          uid,
			})
			println(ret)
			uid = uid + 1
		}
	}
}

//  添加发放碳量记录并且更新用户剩余碳量

func (srv CarbonTransactionService) Create(dto api_types.CreateCarbonTransactionDto) (float64, error) {
	//获取ip地址
	cityCode := ""
	cityInfo, cityErr := baidu.IpToCity(dto.Ip)
	if cityErr == nil {
		cityCode = cityInfo.Content.AddressDetail.Adcode
	}
	//查询场景配置
	scene := srv.repoScene.FindBy(repotypes.CarbonSceneBy{Type: dto.Type})
	if scene.ID == 0 {
		return 0, errors.New("不存在该场景")
	}
	//判断是否有限制
	/*errCheck := NewCarbonTransactionCountLimitService(srv.ctx).CheckLimitAndUpdate(dto.Type, dto.OpenId, scene.MaxCount)
	if errCheck != nil {
		return 0, errCheck
	}*/
	//获取碳量
	carbon := srv.repoScene.GetValue(scene, dto.Value) //增加的碳量
	_, err := NewCarbonService(context.NewMioContext()).IncUserCarbon(srv_types.IncUserCarbonDTO{
		OpenId:       dto.OpenId,
		Type:         dto.Type,
		BizId:        util.UUID(),
		ChangePoint:  carbon,
		AdditionInfo: dto.Info,
		CityCode:     cityCode,
		Uid:          dto.UserId,
	})
	if err != nil {
		return 0, err
	}
	return carbon, nil
}

// Bank 排行榜
func (srv CarbonTransactionService) Bank(dto api_types.GetCarbonTransactionBankDto) ([]api_types.CarbonTransactionBank, int64, error) {
	//1 获取排行榜信息
	todayDate := time.Now().Format("20060102") //当天时间 年月日
	redisKey := fmt.Sprintf(config.RedisKey.UserCarbonRank, todayDate)

	op := &redis.ZRangeBy{
		Max:    "10000000",
		Min:    "0",
		Offset: int64(dto.Offset), // 类似sql的limit, 表示开始偏移量
		Count:  int64(dto.Limit),  // 一次返回多少数据
	}
	bank, err := app.Redis.ZRevRangeByScoreWithScores(contextRedis.Background(), redisKey, op).Result()
	if err != nil {
		return nil, 0, err
	}
	total := app.Redis.ZCard(contextRedis.Background(), redisKey)
	//2 排行榜当前页所有的uid
	uids := make([]int64, 0)
	for _, val := range bank {
		i, errM := strconv.ParseInt(val.Member.(string), 10, 64)
		if errM == nil {
			uids = append(uids, i)
		}
	}
	//3 根据uid获取用户信息
	var userList []entity.User
	if len(uids) != 0 {
		userList, _ = DefaultUserService.GetUserListBy(repository.GetUserListBy{UserIds: uids})
	}
	userMap := make(map[int64]api_types.CarbonUser)
	for _, val := range userList {
		userMap[val.ID] = api_types.CarbonUser{
			Nickname:  val.Nickname,
			AvatarUrl: val.AvatarUrl,
		}
	}
	//4 根据当前页的uid和我的uid去好友表查到存在的uid
	friends, _ := DefaultUserFriendService.GetUserFriendList(3, uids)

	//5 根据用户信息和好友信息进行信息整理
	CarbonBankList := make([]api_types.CarbonTransactionBank, 0)
	for key, val := range bank {
		member, errM := strconv.ParseInt(val.Member.(string), 10, 64)
		if errM != nil {
			continue
		}
		var idFriend = false
		_, ok := friends[member]
		if ok {
			idFriend = true
		} else {
			idFriend = false
		}
		CarbonBankList = append(CarbonBankList, api_types.CarbonTransactionBank{
			Carbon: util.CarbonToRate(val.Score), //分数
			Rank:   int64(key) + 1,               //排行
			User: api_types.CarbonUser{
				AvatarUrl: userMap[member].AvatarUrl,
				Nickname:  userMap[member].Nickname,
				Uid:       member, //用户uid
			},
			IsFriend: idFriend,
		})
	}
	return CarbonBankList, total.Val(), nil
}

// MyBank 我的排名
func (srv CarbonTransactionService) MyBank(dto api_types.GetCarbonTransactionMyBankDto) (api_types.CarbonTransactionMyBank, error) {
	todayDate := time.Now().Format("20060102") //当天时间 年月日
	redisKey := fmt.Sprintf(config.RedisKey.UserCarbonRank, todayDate)
	uidStr := strconv.FormatInt(dto.UserId, 10) //我的uid string

	mySort := app.Redis.ZRank(contextRedis.Background(), redisKey, uidStr)    //我的当前排名
	myCarbon := app.Redis.ZScore(contextRedis.Background(), redisKey, uidStr) //我的碳量
	allCount := app.Redis.ZCard(contextRedis.Background(), redisKey)          //总人数

	mySortDec := decimal.NewFromInt(0)
	if mySort.Val() != 0 {
		mySortDec = decimal.NewFromInt(mySort.Val()).Add(decimal.NewFromInt(1))
	}

	allCountDec := decimal.NewFromInt(allCount.Val())

	//超越**%用户
	var overPer string
	if !allCountDec.IsZero() {
		overPer = mySortDec.Div(allCountDec).Round(2).Mul(decimal.NewFromInt(100)).String() + "%"
	}
	//我的当前排名
	myRank := allCountDec.Sub(mySortDec).Add(decimal.NewFromInt(1))

	//我的信息
	user, err := DefaultUserService.GetUserById(dto.UserId)
	if err != nil {
		return api_types.CarbonTransactionMyBank{}, err
	}
	myBank := api_types.CarbonTransactionMyBank{
		OverPer: overPer,
		Carbon:  util.CarbonToRate(myCarbon.Val()),
		Rank:    myRank.String(),
		User: api_types.CarbonUser{
			AvatarUrl: user.AvatarUrl, //头像
			Nickname:  user.Nickname,
			Uid:       user.ID}, //昵称
	}
	return myBank, nil
}

//GetTodayCarbon  获取今日碳量
func (srv CarbonTransactionService) GetTodayCarbon(uid int64) float64 {
	uidStr := strconv.FormatInt(uid, 10)       //我的uid string
	todayDate := time.Now().Format("20060102") //当天时间 年月日
	redisKey := fmt.Sprintf(config.RedisKey.UserCarbonRank, todayDate)
	myCarbon := app.Redis.ZScore(contextRedis.Background(), redisKey, uidStr) //我的碳量
	return myCarbon.Val()
}

type KVPair struct {
	Key entity.CarbonTransactionType
	Val float64
}

func (srv CarbonTransactionService) Classify(dto api_types.GetCarbonTransactionClassifyDto) (retDto api_types.CarbonTransactionClassify, err error) {
	UserIdString := strconv.FormatInt(dto.UserId, 10) //我的uid string
	dataStr := app.Redis.HGet(contextRedis.Background(), config.RedisKey.UserCarbonClassify, UserIdString)
	if dataStr.Val() == "" {
		//用默认的
		retDto.List = []api_types.CarbonTransactionClassifyList{{Key: entity.CARBON_STEP.Text(), Val: 0}, {Key: entity.CARBON_BIKE_RIDE.Text(), Val: 0}, {Key: entity.CARBON_ECAR.Text(), Val: 0}, {Key: "其他", Val: 0}}
		retDto.Cover = entity.CARBON_STEP.Cover()
		retDto.Total = 0
		return
	}
	var dataMap map[entity.CarbonTransactionType]float64
	err = json.Unmarshal([]byte(dataStr.Val()), &dataMap)
	if err != nil {
		return
	}

	//map转化成切片,方便排序
	tmpList := make([]KVPair, 0)
	for k, v := range dataMap {
		tmpList = append(tmpList, KVPair{Key: k, Val: v})
	}
	//排序
	sort.Slice(tmpList, func(i, j int) bool {
		return tmpList[i].Val < tmpList[j].Val // 升序
	})
	//整理
	ret := make([]api_types.CarbonTransactionClassifyList, 0)
	other := 0.00 //其他碳量
	total := 0.00 //总碳量
	for i, _ := range tmpList {
		n := tmpList[len(tmpList)-1-i]
		total += n.Val
		if i == 0 {
			ret = append(ret, api_types.CarbonTransactionClassifyList{
				Val: n.Val,
				Key: n.Key.Text(),
			})
			retDto.Cover = n.Key.Cover() //只有第一个有封面
		} else if i == 1 || i == 2 {
			ret = append(ret, api_types.CarbonTransactionClassifyList{
				Val: n.Val,
				Key: n.Key.Text(),
			})
		} else {
			other += n.Val
		}
	}
	ret = append(ret, api_types.CarbonTransactionClassifyList{Key: "其他", Val: other})
	retDto.Total = total
	retDto.List = ret
	return
}

// History 我的减碳成就-近2周减碳
func (srv CarbonTransactionService) History(dto api_types.GetCarbonTransactionHistoryDto) ([]entity.CarbonTransactionDay, error) {
	list, err := srv.repoDay.GetList(repotypes.GetCarbonTransactionDayGetListDO{
		StartTime: dto.StartTime,
		EndTime:   dto.EndTime,
		UserId:    dto.UserId,
		OrderBy:   entity.OrderByList{entity.OrderByCarbonTranDayVDate},
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Info 我的减碳成就-基础信息
func (srv CarbonTransactionService) Info(dto api_types.GetCarbonTransactionInfoDto) (api_types.CarbonTransactionInfo, error) {
	user, err := DefaultUserService.GetUserById(dto.UserId)
	if err != nil {
		return api_types.CarbonTransactionInfo{}, err
	}

	carbonInfo, err := NewCarbonService(context.NewMioContext()).FindByUserId(dto.UserId)
	if err != nil {
		return api_types.CarbonTransactionInfo{}, err
	}

	TreeNum, TreeNumMsg := util.CarbonToTree(carbonInfo.Carbon)
	carbonToday := srv.GetTodayCarbon(dto.UserId) //今日碳量
	info := api_types.CarbonTransactionInfo{
		RegisterDayNum: timeutils.Now().GetDiffDays(time.Now(), user.Time.Time),
		Carbon:         util.CarbonToRate(carbonInfo.Carbon),
		CarbonToday:    util.CarbonToRate(carbonToday),
		TreeNum:        TreeNum,
		TreeNumMsg:     TreeNumMsg,
		User: api_types.CarbonUser{
			AvatarUrl: user.AvatarUrl, //头像
			Nickname:  user.Nickname,  //昵称
			Uid:       user.ID},
	}
	return info, nil
}

func (srv CarbonTransactionService) AddClassify() {
	list := srv.repo.GetListBy(repotypes.GetCarbonTransactionListByDO{
		EndTime: time.Now().AddDate(0, 0, 0).Format("2006-01-02"),
	})
	DateMap := make(map[int64]map[entity.CarbonTransactionType]float64)
	for _, by := range list {
		_, ok := DateMap[by.UserId]
		if !ok {
			DateMap[by.UserId] = make(map[entity.CarbonTransactionType]float64)
		}
		DateMap[by.UserId][by.Type] = by.Sum
	}
	for k, v := range DateMap {
		marshal, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("Map转化为byte数组失败,异常:%s\n", err)
			return
		}
		UserIdString := strconv.FormatInt(k, 10) //我的uid string
		app.Redis.HDel(contextRedis.Background(), config.RedisKey.UserCarbonClassify, UserIdString)
		ret := app.Redis.HSet(contextRedis.Background(), config.RedisKey.UserCarbonClassify, UserIdString, string(marshal))
		println(ret)
	}
}

//每天总结碳量
func (srv CarbonTransactionService) AddHistory() {
	list := srv.repo.GetListByDay(repotypes.GetCarbonTransactionListByDO{
		StartTime: time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		EndTime:   time.Now().AddDate(0, 0, 0).Format("2006-01-02"),
	})
	for _, v := range list {
		srv.repoDay.Create(&entity.CarbonTransactionDay{
			OpenId: v.Openid,
			UserId: v.UserId,
			VDate:  time.Now().AddDate(0, 0, -1),
			Value:  v.Sum,
		})
	}
}
