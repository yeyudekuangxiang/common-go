package jhx

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"math/rand"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/httputil"
	"mio/pkg/errno"
	platformUtil "mio/pkg/platform"
	"strconv"
	"strings"
	"time"
)

type jhxOption struct {
	Domain    string
	AppId     string
	Version   string
	Timestamp string
	Nonce     string
}

type Options func(options *jhxOption)

func NewJhxService(ctx *context.MioContext, jhxOptions ...Options) *Service {
	options := &jhxOption{
		Domain:    "http://m.jinhuaxing.com.cn/api",
		AppId:     "2498728d209d",
		Version:   "1.0",
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		Nonce:     strconv.Itoa(rand.Int()),
	}

	for i := range jhxOptions {
		jhxOptions[i](options)
	}

	return &Service{
		ctx:    ctx,
		option: options,
	}
}

type Service struct {
	ctx    *context.MioContext
	option *jhxOption
}

func WithJhxDomain(domain string) Options {
	return func(option *jhxOption) {
		option.Domain = domain
	}
}

func WithJhxAppId(appId string) Options {
	return func(option *jhxOption) {
		option.AppId = appId
	}
}

func WithJhxTimestamp(timestamp string) Options {
	return func(option *jhxOption) {
		option.Timestamp = timestamp
	}
}

func WithJhxNonce(nonce string) Options {
	return func(option *jhxOption) {
		option.Nonce = nonce
	}
}

//发放券码
func (srv Service) SendCoupon(typeId int64, user entity.User) (string, error) {
	tradeNo, err := srv.TicketCreate(typeId, user)
	if err != nil {
		return "", err
	}
	return tradeNo, nil
}

func (srv Service) TicketCreate(typeId int64, user entity.User) (string, error) {
	commonParams := srv.getCommonParams()
	rand.Seed(time.Now().UnixNano())
	tradeNo := "jhx" + strconv.FormatInt(time.Now().UnixMilli(), 10) + strconv.FormatInt(rand.Int63(), 10)
	commonParams["tradeno"] = tradeNo
	commonParams["mobile"] = user.PhoneNumber

	commonParams["sign"] = strings.ToUpper(platformUtil.GetSign(commonParams, "", "&"))

	url := srv.option.Domain + "/busticket/ticket_create"
	body, err := httputil.PostJson(url, commonParams)

	fmt.Printf("ticket_create response body: %s\n", body)
	if err != nil {
		return "", err
	}

	response := jhxCommonResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Unmarshal body: %s\n", err.Error())
		return "", err
	}

	if response.Code != 0 {
		return "", errors.New(response.Msg)
	}

	ticketCreateResponse := &jhxTicketCreateResponse{}
	err = util.MapTo(response.Data, &ticketCreateResponse)
	if err != nil {
		return "", err
	}

	var expireTime time.Time
	if ticketCreateResponse.ExpireTime == "" {
		expireTime = time.Now().AddDate(1, 0, 0)
	} else {
		expireTime, _ = time.Parse("2006-01-02", ticketCreateResponse.ExpireTime)
	}

	_, err = app.RpcService.CouponRpcSrv.SendCoupon(srv.ctx, &couponclient.SendCouponReq{
		CouponCardTypeId:     typeId,
		CouponCardQrcodeText: ticketCreateResponse.QrCodeStr,
		UserId:               user.ID,
		BizId:                tradeNo,
		StartTime:            time.Now().UnixMilli(),
		EndTime:              expireTime.UnixMilli(),
	})
	if err != nil {
		return "", err
	}

	return tradeNo, nil
}

//消费通知
func (srv Service) TicketNotify(sign string, params map[string]interface{}) error {
	if err := platformUtil.CheckSign(sign, params, "d7b47f379109", "&"); err != nil {
		return err
	}

	ticketNotify := TicketNotify{}
	_ = util.MapTo(&params, &ticketNotify)

	//查询库 根据tradeno获取券码
	coupon, err := app.RpcService.CouponRpcSrv.FindCoupon(srv.ctx, &couponclient.FindCouponReq{
		CouponCardTypeId: 1000,
		BizId:            ticketNotify.Tradeno,
	})
	if err != nil {
		return err
	}

	//如果 status 相等 不处理 返回 nil
	j, _ := strconv.ParseInt(ticketNotify.Status, 10, 32)
	status := int32(j)

	if !coupon.Exist {
		return errno.ErrCommon.WithMessage("券码不存在")
	}

	if coupon.CouponInfo.UsedStatus == status {
		return errno.ErrCommon.WithMessage("该券码已失效")
	}

	//如果 status 不等 根据 tradeno 更新status,used_time 返回nil
	_, err = app.RpcService.CouponRpcSrv.UpdateCouponUsedStatus(srv.ctx, &couponclient.UpdateCouponUsedStatusReq{
		CouponCardId: coupon.CouponInfo.CouponCardId,
		UsedStatus:   status,
		UsedTime:     time.Now().UnixMilli(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (srv Service) TicketStatus(tradeno string) (*JhxTicketStatusResponse, error) {
	commonParams := srv.getCommonParams()
	commonParams["tradeno"] = tradeno
	commonParams["sign"] = strings.ToUpper(platformUtil.GetSign(commonParams, "", "&"))

	url := srv.option.Domain + "/busticket/ticket_status"

	body, err := httputil.PostJson(url, commonParams)
	fmt.Printf("%s\n", body)
	if err != nil {
		return &JhxTicketStatusResponse{}, err
	}
	response := jhxCommonResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return &JhxTicketStatusResponse{}, err
	}
	if response.Code != 0 {
		return &JhxTicketStatusResponse{}, errno.ErrCommon.WithMessage(response.Msg)
	}
	ticketStatusResponse := &JhxTicketStatusResponse{}
	err = util.MapTo(response.Data, ticketStatusResponse)
	if err != nil {
		return &JhxTicketStatusResponse{}, err
	}
	//返回状态
	return ticketStatusResponse, nil
}

//绑定回调
func (srv Service) BindSuccess(params map[string]interface{}) error {
	commonParams := srv.getCommonParams()
	commonParams["status"] = params["status"]
	commonParams["mobile"] = params["mobile"]
	commonParams["sign"] = strings.ToUpper(platformUtil.GetSign(commonParams, "", "&"))
	url := srv.option.Domain + "/busticket/bind_success"
	body, err := httputil.PostJson(url, commonParams)
	fmt.Printf("bind_success response body: %s\n", body)
	if err != nil {
		return err
	}
	response := jhxCommonResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("bind_success Unmarshal body: %s\n", err.Error())
		return err
	}
	if response.Code != 0 {
		return errno.ErrCommon.WithMessage(response.Msg)
	}
	return nil
}

//创建气泡数据
func (srv Service) PreCollectPoint(sign string, params map[string]interface{}) error {
	if err := platformUtil.CheckSign(sign, params, "", "&"); err != nil {
		return err

	}

	//根据 platform_member_id 获取 openid
	var openId string
	sceneUser := repository.DefaultBdSceneUserRepository.FindPlatformUserByPlatformUserId(params["memberId"].(string), params["platformKey"].(string))
	if sceneUser.ID != 0 {
		openId = sceneUser.OpenId
	}
	//创建数据
	fromString, err := decimal.NewFromString(params["amount"].(string))
	if err != nil {
		return err
	}

	point := fromString.Mul(decimal.NewFromInt(10)).Round(2).IntPart()
	err = repository.DefaultBdScenePrePointRepository.Create(&entity.BdScenePrePoint{
		PlatformKey:    params["platformKey"].(string),
		PlatformUserId: params["memberId"].(string),
		Point:          point,
		OpenId:         openId,
		Status:         1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

//获取气泡数据
func (srv Service) GetPreCollectPointList(sign string, params map[string]interface{}) ([]entity.BdScenePrePoint, int64, error) {
	if err := platformUtil.CheckSign(sign, params, "", "&"); err != nil {
		return nil, 0, err
	}

	getPreCollect := GetPreCollect{}
	_ = util.MapTo(&params, &getPreCollect)

	var items []entity.BdScenePrePoint
	var point int64

	//获取pre_point数据
	scenePointCondition := repository.GetScenePrePoint{
		PlatformKey: getPreCollect.PlatformKey,
		StartTime:   time.Now().AddDate(0, 0, -7),
		EndTime:     time.Now(),
		Status:      1,
	}
	//根据 platform_member_id 获取 openid
	sceneUserCondition := repository.GetSceneUserOne{
		PlatformKey: getPreCollect.PlatformKey,
	}

	if getPreCollect.MemberId != "" {
		scenePointCondition.PlatformUserId = getPreCollect.MemberId
		sceneUserCondition.PlatformUserId = getPreCollect.MemberId
	}

	if getPreCollect.OpenId != "" {
		scenePointCondition.OpenId = getPreCollect.OpenId
		sceneUserCondition.OpenId = getPreCollect.OpenId
	}

	items, _, err := repository.DefaultBdScenePrePointRepository.FindBy(scenePointCondition)
	if err != nil {
		return items, 0, err
	}

	sceneUser := repository.DefaultBdSceneUserRepository.FindOne(sceneUserCondition)
	if sceneUser.ID != 0 {
		pointInfo := repository.NewPointRepository(srv.ctx).FindBy(repository.FindPointBy{OpenId: sceneUser.OpenId})
		point = pointInfo.Balance
	} else if getPreCollect.OpenId != "" {
		pointInfo := repository.NewPointRepository(srv.ctx).FindBy(repository.FindPointBy{OpenId: getPreCollect.OpenId})
		point = pointInfo.Balance
	}
	return items, point, nil
}

//消费气泡数据
func (srv Service) CollectPoint(sign string, params map[string]interface{}) (int64, error) {
	if err := platformUtil.CheckSign(sign, params, "", "&"); err != nil {
		return 0, err
	}

	collect := Collect{}
	_ = util.MapTo(&params, &collect)

	scene := repository.DefaultBdSceneRepository.FindByCh(collect.PlatformKey)
	if scene.Key == "" || scene.Key == "e" {
		return 0, errno.ErrCommon.WithMessage("渠道查询失败")
	}

	sceneUserCondition := repository.GetSceneUserOne{PlatformKey: collect.PlatformKey}
	if collect.MemberId != "" {
		sceneUserCondition.PlatformUserId = collect.MemberId
	}

	if collect.OpenId != "" {
		sceneUserCondition.OpenId = collect.OpenId
	}

	sceneUser := repository.DefaultBdSceneUserRepository.FindOne(sceneUserCondition)
	if sceneUser.ID == 0 {
		return 0, errno.ErrCommon.WithMessage("未找到绑定关系")
	}

	//获取pre_point数据 one limit
	id, _ := strconv.ParseInt(collect.PrePointId, 10, 64)
	one, err := repository.DefaultBdScenePrePointRepository.FindOne(repository.GetScenePrePoint{
		PlatformKey:    sceneUser.PlatformKey,
		PlatformUserId: sceneUser.PlatformUserId,
		Id:             id,
		Status:         1,
	})
	if err != nil {
		return 0, errno.ErrRecordNotFound
	}

	//检查上限
	var isHalf bool
	var halfPoint int64
	timeStr := time.Now().Format("2006-01-02")
	key := timeStr + ":prePoint:" + scene.Ch + sceneUser.PlatformUserId + sceneUser.Phone
	lastPoint, _ := strconv.ParseInt(app.Redis.Get(srv.ctx, key).Val(), 10, 64)
	incPoint := one.Point
	totalPoint := lastPoint + incPoint
	if lastPoint >= int64(scene.PrePointLimit) {
		return 0, errno.ErrCommon.WithMessage("今日获取积分已达到上限")
	}

	if totalPoint > int64(scene.PrePointLimit) {
		p := incPoint
		incPoint = int64(scene.PrePointLimit) - lastPoint
		totalPoint = int64(scene.PrePointLimit)
		isHalf = true
		halfPoint = p - incPoint
	}

	app.Redis.Set(srv.ctx, key, totalPoint, 24*time.Hour)
	//积分
	point, err := service.NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:      sceneUser.OpenId,
		Type:        entity.POINT_JHX,
		BizId:       util.UUID(),
		ChangePoint: incPoint,
		AdminId:     0,
		Note:        collect.PlatformKey + "#" + one.Tradeno,
	})
	if err != nil {
		return 0, err
	}

	//更新pre_point对应数据
	one.Status = 2
	one.UpdatedAt = time.Now()
	if isHalf {
		one.Status = 1
		one.Point = halfPoint
	}

	err = repository.DefaultBdScenePrePointRepository.Save(&one)
	if err != nil {
		return 0, err
	}

	return point, nil
}

func (srv Service) getCommonParams() map[string]interface{} {
	params := make(map[string]interface{}, 0)
	params["version"] = srv.option.Version
	params["appid"] = srv.option.AppId
	params["timestamp"] = srv.option.Timestamp
	params["nonce"] = srv.option.Nonce
	return params
}
