package activity

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/commontool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"
	"gitlab.miotech.com/miotech-application/backend/common-go/wxapp"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/activity/activityclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/carbonpk/carbonpk"
	"mio/config"
	"mio/internal/pkg/core/app"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/duiba"
	"mio/pkg/errno"
	"strconv"
	"strings"
	"time"
)

var (
	ZeroActivityStartTime, _ = time.Parse("2006-01-02 15:04:05", "2022-04-13 00:00:00")
)
var DefaultZeroService = ZeroService{
	repo:       repository.NewPointTransactionRepository(mioContext.NewMioContext()),
	repoOrder:  repository.NewOrderRepository(app.DB),
	repoInvite: repository.NewInviteRepository(mioContext.NewMioContext()),
}

type ZeroService struct {
	repo       *repository.PointTransactionRepository
	repoOrder  repository.OrderRepository
	repoInvite *repository.InviteRepository
}

func NewZeroService(ctx *mioContext.MioContext) *ZeroService {
	return &ZeroService{
		repo:      repository.NewPointTransactionRepository(ctx),
		repoOrder: repository.NewOrderRepository(app.DB),
	}
}

func (srv ZeroService) AutoLogin(userId int64, short string) (string, error) {
	userInfo, err := service.DefaultUserService.GetUserById(userId)
	if err != nil {
		return "", err
	}
	if userInfo.ID == 0 {
		return "", errno.ErrCommon.WithMessage("未查询到用户信息")
	}

	//此方法只能使用一次
	isNewUser, err := srv.IsNewUser(userId, userInfo.Time.Time)
	if err != nil {
		return "", err
	}
	path := "https://88543.activity-12.m.duiba.com.cn/aaw/haggle/index?opId=195380425492081&dbnewopen&newChannelType=3"
	if short != "" {
		p, err := srv.GetUrlByShort(short)
		if err != nil {
			app.Logger.Error(userId, short, err)
		}
		if p != "" {
			path = p
		}
	}

	return duiba.DefaultDuiBaService.AutoLoginOpenId(service.AutoLoginOpenIdParam{
		UserId:  userId,
		OpenId:  userInfo.OpenId,
		Path:    path,
		DCustom: fmt.Sprintf("avatar=%s&nickname=%s&newUser=%d", userInfo.AvatarUrl, userInfo.Nickname, isNewUser),
	})
}
func (srv ZeroService) StoreUrl(url string) (string, error) {
	key := encrypttool.Md5(url)
	redisKey := fmt.Sprintf(config.RedisKey.DuiBaShortUrl, key)
	err := app.Redis.Set(context.Background(), redisKey, url, time.Hour*10*24).Err()
	if err != nil {
		return "", err
	}
	return key, nil
}
func (srv ZeroService) GetUrlByShort(short string) (string, error) {
	redisKey := fmt.Sprintf(config.RedisKey.DuiBaShortUrl, short)
	u, err := app.Redis.Get(context.Background(), redisKey).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return u, nil
}
func (srv ZeroService) IsNewUser(userId int64, createTime time.Time) (int, error) {
	//用户创建时间在活动开始时间之前
	if createTime.Before(ZeroActivityStartTime) {
		return 0, nil
	}
	return 1, nil

	//兑吧自己有判断 不用我们记录了
	//判断是否已经助力过
	/*redisKey := fmt.Sprintf(config.RedisKey.ActivityZeroIsNewUser, userId)
	result, err := app.Redis.SetNX(context.Background(), redisKey, time.Now().Unix(), time.Hour*24*30).Result()
	if err != nil && err != redis.Nil {
		return 0, err
	}
	return commontool.Ternary(result, 1, 0).Int(), nil*/
}

type DuiBaActivity struct {
	ActivityId  string
	ActivityUrl string
	StartTime   time.Time
	EndTime     time.Time
}

const DUIBAIndex = "https://88543.activity-12.m.duiba.com.cn/chw/visual-editor/skins?id=239935"

func (srv ZeroService) DuiBaStoreUrl(activityId string, url string) (string, error) {
	key := activityId + "_" + encrypttool.Md5(url)
	redisKey := fmt.Sprintf(config.RedisKey.DuiBaShortUrl, key)
	err := app.Redis.Set(context.Background(), redisKey, url, time.Hour*10*24).Err()
	if err != nil {
		return "", err
	}
	return key, nil
}
func (srv ZeroService) GetDuiBaUrlByShort(short string) (string, error) {
	redisKey := fmt.Sprintf(config.RedisKey.DuiBaShortUrl, short)
	u, err := app.Redis.Get(context.Background(), redisKey).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return u, nil
}
func (srv ZeroService) DuiBaAutoLogin(userId int64, activityId, short, thirdParty string, cip string) (string, error) {
	//userId = 30
	userInfo, err := service.DefaultUserService.GetUserById(userId)
	if err != nil {
		return "", err
	}
	if userInfo.ID == 0 {
		return "", errno.ErrUserNotFound
	}

	path := DUIBAIndex
	isNewUser := false

	pointService := service.NewDuiBaActivityService(mioContext.NewMioContext())
	activity, err := pointService.FindActivity(activityId)
	if err != nil {
		return "", err
	}

	if activity.ID != 0 {
		isNewUser, err = srv.IsDuiBaActivityNewUser(activityId, userId)
		if err != nil {
			return "", err
		}

		path = activity.ActivityUrl

		if short != "" {
			p, err := srv.GetDuiBaUrlByShort(short)
			if err != nil {
				app.Logger.Error(userId, short, err)
			}
			if p != "" {
				path = p
			}
		}
	}

	//检测用户风险等级
	userRiskRankParam := wxapp.UserRiskRankParam{
		OpenId:   userInfo.OpenId,
		Scene:    0,
		ClientIp: cip,
	}
	//风险等级4,代表不需要验证.0-3代表需要做出限制
	if activity.RiskLimit < 4 {
		//判断用户手机号,警告:必须是非首页
		if userInfo.PhoneNumber == "" {
			return "", errno.ErrMisMatchCondition
		}
		userRiskRankParam.MobileNo = userInfo.PhoneNumber
	}
	/*resp, err := service.DefaultUserService.CheckUserRisk(userRiskRankParam)
	if err != nil && activityId != "index" {
		app.Logger.Info("DuiBaAutoLogin 风险等级查询查询出错", err.Error())
	} else {
		if resp.RiskRank > activity.RiskLimit {
			return "", nil
			//return "", errors.New("该活动仅限部分地区用户参与,请看看其他活动")
		}
	}*/
	if userInfo.Risk > activity.RiskLimit {
		return "", nil
		//return "", errors.New("该活动仅限部分地区用户参与,请看看其他活动")
	}

	vip := 0
	/*if thirdParty == "thirdParty" && isNewUser {
		vip = 2
	}*/

	switch activity.VipType {
	case entity.DuiBaActivityInviteActivity:
		//jira: https://jira.miotech.com/browse/MP2C-1681?goToView=5
		if userInfo.Risk == 4 {
			break
		}
		inviteCount, err := srv.repoInvite.GetInviteRewardFenQun(repotypes.GetInviteTotalDO{
			StartTime: "2022-09-15:00:00:01",
			EndTime:   "2022-09-22:00:00:01",
			Openid:    userInfo.OpenId,
		})
		if err != nil {
			break
		}
		switch {
		case inviteCount >= 3 && inviteCount <= 5:
			{
				vip = 52
				break
			}
		case inviteCount >= 6 && inviteCount <= 14:
			{
				vip = 53
				break
			}
		case inviteCount >= 15:
			{
				vip = 54
				orderTotal := srv.repoOrder.GetOrderTotalByItemId(repotypes.GetOrderTotalByItemIdDO{
					Openid:    userInfo.OpenId,
					StartTime: "2022-09-15:00:00:01",
					EndTime:   "2022-09-22:00:00:01"})
				if orderTotal >= 1 && inviteCount >= 21 {
					vip = 55
				}
				break
			}
		default:
			{
				break
			}
		}
		break

	case entity.DuiBaActivityVipTypeNewUser:
		if thirdParty == "thirdParty" && isNewUser {
			vip = 2
		}
		break
	case entity.DuiBaActivityRecyclingPublicWelfareWeekActivity:
		//需求地址：https://jira.miotech.com/browse/MP2C-1591
		var pointTypes = []string{"RECYCLING_CLOTHING", "RECYCLING_COMPUTER", "RECYCLING_APPLIANCE", "RECYCLING_BOOK"}
		count := srv.repo.GetListByFenQunCount(repository.GetPointTransactionListByQun{
			StartTime: "2022-09-07:00:00:01",
			EndTime:   "2022-09-14:00:00:01",
			Types:     pointTypes,
			OpenId:    userInfo.OpenId,
		})
		var ItemIdSlice = []string{"cbddf0af60f402f717b0987b79709209", "b00064a760f400a42850b68e1f783c22"}
		orderTotal := srv.repoOrder.GetOrderTotalByItemId(repotypes.GetOrderTotalByItemIdDO{
			Openid:      userInfo.OpenId,
			ItemIdSlice: ItemIdSlice,
			StartTime:   "2022-09-07:00:00:01",
			EndTime:     "2022-09-14:00:00:01"})
		if count >= 1 && orderTotal >= 1 {
			vip = 56
		}
	case entity.DuiBaActivityIsPhoneAnniversaryActivity:
		//需求地址：https://confluence.miotech.com/pages/viewpage.action?pageId=26613756
		var pointTypes = []string{"STEP", "COFFEE_CUP", "BIKE_RIDE", "ECAR"}
		count := srv.repo.GetListByFenQunCount(repository.GetPointTransactionListByQun{
			StartTime: "2022-08-01:00:00:01",
			EndTime:   "2022-08-15:00:00:01",
			Types:     pointTypes,
			OpenId:    userInfo.OpenId,
		})
		switch {
		case count == 3:
			{
				vip = 57
			}
		case count >= 4 && count <= 7:
			{
				vip = 58
			}
		case count >= 8 && count <= 14:
			{
				vip = 59
			}
		default:
			{
				break
			}
		}
		break

	case entity.DuiBaActivityYtxLuckyDrawActivity:
		if userInfo.PhoneNumber == "" {
			return "", errno.ErrMisMatchCondition
		}
		//获取详情
		user, err := app.RpcService.ActivityRpcSrv.DetailActivityThirdUser(context.Background(), &activityclient.DetailActivityThirdUserReq{
			ActivityId: 3,
			Phone:      userInfo.PhoneNumber,
		})
		if err != nil {
			return "", err
		}
		if !user.Exist {
			return "", errno.ErrMisMatchCondition
		}
		break

	case entity.DuiBaActivityCaronPkLuckyDrawActivity:
		cid := strconv.FormatInt(activity.Cid, 10)
		log, err := app.RpcService.CarbonPkRpcSrv.DetailCarbonPkRewardLog(mioContext.NewMioContext(), &carbonpk.DetailCarbonPkRewardReq{
			RewardUrl: "/pages/duiba_v2/duiba/index?activityId=" + activityId + "&cid=" + cid + "&bind=true&hideShare=false",
			UserId:    userId,
		})
		if err != nil {
			return "", errno.ErrMisMatchCondition
		}
		if !log.Exist {
			return "", errno.ErrMisMatchCondition
		}
		break

	case entity.DuiBaActivityRedBlackRankActivity:
		isExit := false
		VipOpenidArr := strings.Split(activity.VipOpenid, ",")
		for _, s := range VipOpenidArr {
			if s == userInfo.OpenId {
				isExit = true
				break
			}
		}
		if !isExit {
			return "", errno.ErrMisMatchCondition
		}
		vip = activity.VipId
		break
	case entity.DuiBaActivityCrazyThursdayActivity:
		CrazyThursdayActivityId := strings.Split(activityId, "_")
		startTime := timetool.MustParse("2006-01-02", CrazyThursdayActivityId[1])
		if startTime.IsZero() {
			return "", errno.ErrMisMatchCondition
		}
		endTime := timetool.MustParse("2006-01-02", CrazyThursdayActivityId[2])
		if endTime.IsZero() {
			return "", errno.ErrMisMatchCondition
		}
		var crazyThursdayActivitySql = `select count(*)  from "duiba_order" where order_id in (
select order_id from order_item where item_id in (select  product_item_id  from product_item where title like '%疯狂星期四%' 
)   order  by id desc)  and finish_time > ? and finish_time < ? and user_id = ?`
		var crazyThursdayActivityTotal int64
		err = app.DB.Raw(crazyThursdayActivitySql, startTime.UnixMilli(), endTime.UnixMilli(), userId).Scan(&crazyThursdayActivityTotal).Error
		if err != nil {
			app.Logger.Error("查询疯狂星期四订单失败", err)
			return "", errno.ErrMisMatchCondition
		}
		if crazyThursdayActivityTotal <= 0 {
			return "", errno.ErrMisMatchCondition
		}
		vip = activity.VipId
		break
	case entity.DuiBaActivityStarChargeActivity:
		rdsKey := fmt.Sprintf("%s:%s", config.RedisKey.StarCharge, "starCharge-luckyDraw")
		if !app.Redis.SIsMember(context.Background(), rdsKey, userInfo.ID).Val() {
			return "", errno.ErrMisMatchCondition
		}
		vip = activity.VipId
		break
	case entity.DuiBaActivityCarbon618Activity:
		carbon618ActivityId := strings.Split(activityId, "_")
		startTime := timetool.MustParse("2006-01-02", carbon618ActivityId[1])
		if startTime.IsZero() {
			return "", errno.ErrMisMatchCondition
		}
		endTime := timetool.MustParse("2006-01-02", carbon618ActivityId[2])
		if endTime.IsZero() {
			return "", errno.ErrMisMatchCondition
		}
		var carbon618ActivitySql = `select count(*)  from "duiba_order" where order_id in (
select order_id from order_item where item_id in (select  product_item_id  from product_item where title like '低碳618%' 
)   order  by id desc)  and finish_time > ? and finish_time < ? and user_id = ?`
		var carbon618ActivityTotal int64
		err = app.DB.Raw(carbon618ActivitySql, startTime.UnixMilli(), endTime.UnixMilli(), userId).Scan(&carbon618ActivityTotal).Error
		if err != nil {
			app.Logger.Error("查询疯狂星期四订单失败", err)
			return "", errno.ErrMisMatchCondition
		}
		if carbon618ActivityTotal <= 0 {
			return "", errno.ErrMisMatchCondition
		}
		vip = activity.VipId
	default:
		break
	}

	//如果是签到的活动，对该用户签到提醒延期处理，以防签到完，24小时之内，对用户再次提醒造成打扰
	/*if activityId == "db_bd_704_qiandao" {
		messageService := messageSrv.MessageService{}
		messageService.ExtensionSignTime(userInfo.OpenId)
	}*/

	isNewUserInt := commontool.Ternary(isNewUser, 1, 0).Int()
	return duiba.DefaultDuiBaService.AutoLoginOpenId(service.AutoLoginOpenIdParam{
		UserId:  userId,
		OpenId:  userInfo.OpenId,
		Path:    path,
		Vip:     vip,
		DCustom: fmt.Sprintf("avatar=%s&nickname=%s&newUser=%d", userInfo.AvatarUrl, userInfo.Nickname, isNewUserInt),
	})
}
func (srv ZeroService) IsDuiBaActivityNewUser(activityId string, userId int64) (bool, error) {
	userInfo, err := service.DefaultUserService.GetUserById(userId)
	if err != nil {
		return false, err
	}
	if userInfo.ID == 0 {
		return false, errno.ErrUserNotFound
	}
	return !userInfo.Time.IsZero() && userInfo.Time.Add(time.Hour*24).After(time.Now()), nil
}
