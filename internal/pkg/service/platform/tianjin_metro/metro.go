package tianjin_metro

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
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

func NewTianjinMetroService(ctx *context.MioContext) *Service {
	return &Service{
		ctx: ctx,
	}
}

type Service struct {
	ctx *context.MioContext
}

func (srv *Service) SendCoupon(typeId int64, amount float64, user entity.User) (string, error) {
	sceneUser := repository.DefaultBdSceneUserRepository.FindOne(repository.GetSceneUserOne{
		PlatformKey: "tianjinmetro",
		OpenId:      user.OpenId,
	})

	if sceneUser.PlatformUserId == "" {
		app.Logger.Errorf("天津地铁 未找到绑定关系 : %s", user.OpenId)
		return "", errno.ErrBindRecordNotFound
	}
	/*
		rand.Seed(time.Now().UnixNano())
		grantV2Request := GrantV3Request{
			AppId:     "264735a59163453d9772f92e1f703123",
			AppSecret: srv.getAppSecret(),
			Ts:        strconv.FormatInt(time.Now().Unix(), 10),
			ReqData: GrantV3ReqData{
				OrderNo:  "tianjinmetro" + util.UUID(),
				PoolCode: srv.option.PoolCode,
				Amount:   amount,
				OpenId:   sceneUser.PlatformUserId,
				Remark:   "lvmiao" + strconv.FormatFloat(amount, 'f', -1, 64) + "元红包",
			},
		}

		url := "https://app.trtpazyz.com/tj-metro-api/open-forward/api/eTicket/allot"
		body, err := httputil.PostJson(url, grantV2Request)
		app.Logger.Infof("天津地铁 grantV2 返回 : %s", body)
		if err != nil {
			return "", err
		}

		response := GrantV3Response{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			app.Logger.Errorf("天津地铁 grantV2 json_decode_err: %s", err.Error())
			return "", err
		}

		if response.SubCode != "0000" {
			app.Logger.Errorf("天津地铁 grantV2 err: %s\n", response.SubMessage)
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
			app.Logger.Errorf("天津地铁 券包 发放错误 : %s\n", err.Error())
			return "", err
		}
	*/
	//	return response.SubData.OrderNo, nil

	return "", nil
}

/*
func (srv *Service) getAppSecret() string {
	t := time.Now().Unix()
	return encrypt.Md5(srv.option.AppId + encrypt.Md5(srv.option.Secret) + strconv.FormatInt(t, 10))
}*/
