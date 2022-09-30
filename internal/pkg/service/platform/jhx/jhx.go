package jhx

import (
	"encoding/json"
	"errors"
	"fmt"
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
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	"mio/pkg/errno"
	"mio/pkg/platform"
	"sort"
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

type JhxOptions func(options *jhxOption)

func NewJhxService(ctx *context.MioContext, jhxOptions ...JhxOptions) *JhxService {
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

	return &JhxService{
		ctx:    ctx,
		option: options,
	}
}

type JhxService struct {
	ctx    *context.MioContext
	option *jhxOption
}

func WithJhxDomain(domain string) JhxOptions {
	return func(option *jhxOption) {
		option.Domain = domain
	}
}

func WithJhxAppId(appId string) JhxOptions {
	return func(option *jhxOption) {
		option.AppId = appId
	}
}

func WithJhxTimestamp(timestamp string) JhxOptions {
	return func(option *jhxOption) {
		option.Timestamp = timestamp
	}
}

func WithJhxNonce(nonce string) JhxOptions {
	return func(option *jhxOption) {
		option.Nonce = nonce
	}
}

func (srv JhxService) TicketCreate(tradeno string, typeId int64, starTime, endTime time.Time, user entity.User) error {
	params := srv.getCommonParams()
	params["tradeno"] = tradeno
	params["mobile"] = user.PhoneNumber
	sign := srv.getSign(params)
	params["sign"] = strings.ToUpper(sign)
	url := srv.option.Domain + "/busticket/ticket_create"
	body, err := httputil.PostJson(url, params)
	fmt.Printf("response body: %s\n", body)
	if err != nil {
		return err
	}
	response := jhxCommonResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Unmarshal body: %s\n", err.Error())
		return err
	}
	if response.Code != 0 {
		return errors.New(response.Msg)
	}
	ticketCreateResponse := &jhxTicketCreateResponse{}
	err = util.MapTo(response.Data, &ticketCreateResponse)
	if err != nil {
		return err
	}
	coupon, err := app.RpcService.CouponRpcSrv.SendCoupon(srv.ctx, &couponclient.SendCouponReq{
		CouponCardTypeId:     typeId,
		CouponCardQrcodeText: ticketCreateResponse.QrCodeStr,
		UserId:               user.ID,
		BizId:                tradeno,
		StartTime:            starTime.UnixMilli(),
		EndTime:              endTime.UnixMilli(),
	})
	fmt.Printf("coupon: %v\n", coupon)
	if err != nil {
		return err
	}
	return nil
}

//消费通知
func (srv JhxService) TicketNotify(sign string, params map[string]interface{}) error {
	scene := service.DefaultBdSceneService.FindByCh("jinhuaxing")
	if scene.Key == "" || scene.Key == "e" {
		return errors.New("渠道查询失败")
	}
	if err := platform.CheckSign(sign, params, scene.Key, "&"); err != nil {
		return err
	}
	//查询库 根据tradeno获取券码
	coupon, err := app.RpcService.CouponRpcSrv.FindCoupon(srv.ctx, &couponclient.FindCouponReq{
		CouponCardTypeId: 123,
		BizId:            params["tradeno"].(string),
	})
	if err != nil {
		return err
	}
	//如果 status 相等 不处理 返回 nil
	j, _ := strconv.ParseInt(params["status"].(string), 10, 32)
	status := int32(j)
	if !coupon.Exist {
		return errors.New("券码不存在")
	}
	if coupon.CouponInfo.UsedStatus == status {
		return errors.New("该券码已失效")
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

func (srv JhxService) TicketStatus(tradeno string) (*jhxTicketStatusResponse, error) {
	params := srv.getCommonParams()
	params["tradeno"] = tradeno
	sign := srv.getSign(params)
	params["sign"] = strings.ToUpper(sign)
	url := srv.option.Domain + "/busticket/ticket_status"
	marshal, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n%s\n", marshal, marshal)
	body, err := httputil.PostJson(url, params)
	fmt.Printf("%s\n", body)
	if err != nil {
		return &jhxTicketStatusResponse{}, err
	}
	response := jhxCommonResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return &jhxTicketStatusResponse{}, err
	}
	if response.Code != 0 {
		return &jhxTicketStatusResponse{}, errors.New(response.Msg)
	}
	ticketStatusResponse := &jhxTicketStatusResponse{}
	err = util.MapTo(response.Data, ticketStatusResponse)
	if err != nil {
		return &jhxTicketStatusResponse{}, err
	}
	//返回状态
	return ticketStatusResponse, nil
}

//创建气泡数据
func (srv JhxService) PreCollectPoint(sign string, params map[string]string) error {
	if err := srv.checkSign(sign, params); err != nil {
		return err
	}
	//根据 platform_member_id 获取 openid
	sceneUser := repository.DefaultBdSceneUserRepository.FindPlatformUserByPlatformUserId(params["memberId"], params["platformKey"])
	if sceneUser.ID == 0 {
		return errors.New("未找到绑定关系")
	}
	//创建数据
	fromString, err := decimal.NewFromString(params["amount"])
	if err != nil {
		return err
	}
	point := fromString.Mul(decimal.NewFromInt(10)).Round(2).String()
	err = repository.DefaultBdScenePrePointRepository.Create(&entity.BdScenePrePoint{
		PlatformKey:    sceneUser.PlatformKey,
		PlatformUserId: sceneUser.PlatformUserId,
		Point:          point,
		OpenId:         sceneUser.OpenId,
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
func (srv JhxService) GetPreCollectPointList(sign string, params map[string]string) ([]entity.BdScenePrePoint, int64, error) {
	if err := srv.checkSign(sign, params); err != nil {
		return nil, 0, err
	}
	//根据 platform_member_id 获取 openid
	sceneUser := repository.DefaultBdSceneUserRepository.FindPlatformUserByPlatformUserId(params["memberId"], params["platformKey"])
	if sceneUser.ID == 0 {
		return nil, 0, errors.New("未找到绑定关系")
	}
	var items []entity.BdScenePrePoint
	var point int64
	//获取pre_point数据
	items, _, err := repository.DefaultBdScenePrePointRepository.FindBy(repository.GetScenePrePoint{
		PlatformKey:    sceneUser.PlatformKey,
		PlatformUserId: sceneUser.PlatformUserId,
		StartTime:      time.Now().AddDate(0, 0, -7).UTC().Format("2006-01-02 15:04:05"),
		EndTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Status:         1,
	})
	if err != nil {
		return items, 0, err
	}
	//获取现有积分
	pointInfo := repository.NewPointRepository(srv.ctx).FindBy(repository.FindPointBy{OpenId: sceneUser.OpenId})
	point = pointInfo.Balance
	return items, point, nil
}

//消费气泡数据
func (srv JhxService) CollectPoint(sign string, params map[string]string) (int64, error) {
	if err := srv.checkSign(sign, params); err != nil {
		return 0, err
	}
	//根据 platform_member_id 获取 openid
	sceneUser := repository.DefaultBdSceneUserRepository.FindPlatformUserByPlatformUserId(params["memberId"], params["platformKey"])
	if sceneUser.ID == 0 {
		return 0, errors.New("未找到绑定关系")
	}
	scene := repository.DefaultBdSceneRepository.FindByCh(params["platformKey"])
	if scene.Key == "" || scene.Key == "e" {
		return 0, errors.New("渠道查询失败")
	}
	//获取pre_point数据 one limit
	id, _ := strconv.ParseInt(params["prePointId"], 10, 64)
	one, err := repository.DefaultBdScenePrePointRepository.FindOne(repository.GetScenePrePoint{
		PlatformKey:    sceneUser.PlatformKey,
		PlatformUserId: sceneUser.PlatformUserId,
		Id:             id,
		Status:         1,
	})
	if err != nil {
		return 0, errno.ErrRecordNotFound
	}
	//调用point_trans.incPoint

	incPoint, _ := strconv.ParseInt(one.Point, 10, 64)
	point, err := service.NewPointService(srv.ctx).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:      sceneUser.OpenId,
		Type:        entity.POINT_JHX,
		BizId:       util.UUID(),
		ChangePoint: incPoint,
		AdminId:     0,
	})
	if err != nil {
		return 0, err
	}
	//更新pre_point对应数据
	one.Status = 2
	one.UpdatedAt = time.Now()
	err = repository.DefaultBdScenePrePointRepository.Save(&one)
	if err != nil {
		return 0, err
	}
	return point, nil
}

func (srv JhxService) MyAccountInfo(sign string, params map[string]string) (*service.UserAccountInfo, error) {
	err := srv.checkSign(sign, params)
	if err != nil {
		return &service.UserAccountInfo{}, err
	}
	//根据 platform_member_id 获取 openid
	sceneUser := repository.DefaultBdSceneUserRepository.FindPlatformUserByPlatformUserId(params["memberId"], params["platformKey"])
	if sceneUser.ID == 0 {
		return &service.UserAccountInfo{}, errors.New("未找到绑定关系")
	}
	userInfo, err := service.DefaultUserService.GetUserByOpenId(sceneUser.OpenId)
	if err != nil {
		return &service.UserAccountInfo{}, err
	}
	accountInfo, err := service.DefaultUserService.AccountInfo(userInfo.ID)
	if err != nil {
		return &service.UserAccountInfo{}, err
	}
	return accountInfo, nil
}

func (srv JhxService) checkSign(sign string, params map[string]string) error {
	md5Sign := srv.getSign(params)
	if sign != md5Sign {
		return errors.New("验签失败")
	}
	return nil
}

// GetSign 签名
func (srv JhxService) getSign(params map[string]string) string {
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += v + "=" + params[v] + "&"
	}
	signStr = strings.TrimRight(signStr, "&")
	return encrypt.Md5(signStr)
}

func (srv JhxService) getCommonParams() map[string]string {
	params := make(map[string]string, 0)
	params["version"] = srv.option.Version
	params["appid"] = srv.option.AppId
	params["timestamp"] = srv.option.Timestamp
	params["nonce"] = srv.option.Nonce
	return params
}
