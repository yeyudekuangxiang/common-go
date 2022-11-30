package ytx

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"math/rand"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	entityActivity "mio/internal/pkg/model/entity/activity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/activity"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	platformUtil "mio/internal/pkg/util/platform"
	"mio/pkg/errno"
	"strconv"
	"time"
)

type ytxOption struct {
	Domain   string
	Secret   string
	PoolCode string
	AppId    string
}

type Options func(options *ytxOption)

//openid:  CpziorTGUL02NrrBqsbbhsAN0Ve4ZMSpPEmgBPAGZOY=
//secret:   a123456
//appid: cc5dec82209c45888620eabec3a29b50
//poolCode: RP202110251300002

func NewYtxService(ctx *context.MioContext, ytxOptions ...Options) *Service {
	options := &ytxOption{
		Domain: "https://apift.ruubypay.com",
	}

	for i := range ytxOptions {
		ytxOptions[i](options)
	}

	return &Service{
		ctx:    ctx,
		option: options,
	}
}

type Service struct {
	ctx    *context.MioContext
	option *ytxOption
}

func WithDomain(domain string) Options {
	return func(option *ytxOption) {
		option.Domain = domain
	}
}

func WithSecret(secret string) Options {
	return func(option *ytxOption) {
		option.Secret = secret
	}
}

func WithPoolCode(poolCode string) Options {
	return func(option *ytxOption) {
		option.PoolCode = poolCode
	}
}

func WithAppId(appId string) Options {
	return func(option *ytxOption) {
		option.AppId = appId
	}
}

//绑定回调
func (srv *Service) BindSuccess(params map[string]interface{}) error {
	synchroRequest := SynchroRequest{
		OpenId:         params["memberId"].(string),
		RegDate:        time.Now().Format("20060102150405"),
		PlatformUserId: params["openId"].(string),
		Ts:             time.Now().UnixMilli(),
	}

	requestParams := make(map[string]interface{}, 0)
	err := util.MapTo(&synchroRequest, &requestParams)

	if err != nil {
		return err
	}

	requestParams["secret"] = srv.option.Secret
	synchroRequest.Signature = platformUtil.GetSign(requestParams, "", "&")

	url := srv.option.Domain + "/markting_activity/network/lvmiao/synchro"
	body, err := httputil.PostJson(url, synchroRequest)
	app.Logger.Infof("亿通行 注册回调 返回 : %s", body)
	if err != nil {
		return err
	}

	response := synchroResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		app.Logger.Errorf("亿通行 json_decode_err: %s", err.Error())
		return err
	}

	if response.ResCode != "0000" {
		app.Logger.Errorf("亿通行 注册回调 错误 : %s", response.ResMessage)
		return errors.New(response.ResMessage)
	}

	return nil
}

func (srv *Service) SendCoupon(typeId int64, amount float64, user entity.User) (string, error) {
	sceneUser := repository.DefaultBdSceneUserRepository.FindOne(repository.GetSceneUserOne{
		PlatformKey: "yitongxing",
		OpenId:      user.OpenId,
	})

	if sceneUser.PlatformUserId == "" {
		app.Logger.Errorf("亿通行 未找到绑定关系 : %s", user.OpenId)
		return "", errno.ErrBindRecordNotFound
	}

	rand.Seed(time.Now().UnixNano())
	grantV2Request := GrantV2Request{
		AppId:     srv.option.AppId,
		AppSecret: srv.getAppSecret(),
		Ts:        strconv.FormatInt(time.Now().Unix(), 10),
		ReqData: GrantV2ReqData{
			OrderNo:  "ytx" + util.UUID(),
			PoolCode: srv.option.PoolCode,
			Amount:   amount,
			OpenId:   sceneUser.PlatformUserId,
			Remark:   "lvmiao" + strconv.FormatFloat(amount, 'f', -1, 64) + "元红包",
		},
	}

	url := srv.option.Domain + "/markting_redenvelopegateway/redenvelope/grantV2"
	body, err := httputil.PostJson(url, grantV2Request)
	app.Logger.Infof("亿通行 grantV2 返回 : %s", body)
	if err != nil {
		return "", err
	}

	response := GrantV2Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		app.Logger.Errorf("亿通行 grantV2 json_decode_err: %s", err.Error())
		return "", err
	}

	respData, _ := json.Marshal(response)

	err = activity.NewYtxLogRepository(srv.ctx).Save(&entityActivity.YtxLog{
		OrderNo:        response.SubData.OrderNo,
		OpenId:         sceneUser.OpenId,
		PlatformUserId: sceneUser.PlatformUserId,
		Amount:         amount,
		AdditionalInfo: string(respData),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	if err != nil {
		app.Logger.Errorf("亿通行 save_log_err : %s", err.Error())
	}

	if response.SubCode != "0000" {
		app.Logger.Errorf("亿通行 grantV2 err: %s\n", response.SubMessage)
		return "", errors.New(response.SubMessage)
	}

	//记录
	_, err = app.RpcService.CouponRpcSrv.SendCoupon(srv.ctx, &couponclient.SendCouponReq{
		CouponCardTypeId: typeId,
		UserId:           user.ID,
		BizId:            response.SubData.OrderNo,
		CouponCardTitle:  "亿通行" + fmt.Sprintf("%.0f", amount) + "元出行红包",
		StartTime:        time.Now().UnixMilli(),
		EndTime:          time.Now().AddDate(0, 0, 90).UnixMilli(),
	})

	if err != nil {
		app.Logger.Errorf("亿通行 券包 发放错误 : %s\n", err.Error())
		return "", err
	}

	return response.SubData.OrderNo, nil
}

func (srv *Service) getAppSecret() string {
	t := time.Now().Unix()
	return encrypt.Md5(srv.option.AppId + encrypt.Md5(srv.option.Secret) + strconv.FormatInt(t, 10))
}