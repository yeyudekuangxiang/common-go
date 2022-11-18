package tianjin_metro

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	"mio/pkg/errno"
	"strconv"
	"strings"
	"time"
)

func NewTianjinMetroService(ctx *context.MioContext) *Service {
	return &Service{
		ctx: ctx,
	}
}

type Service struct {
	ctx *context.MioContext
}

func (srv *Service) SendCoupon(typeId int64, amount float64, user entity.User) (string, error) {
	//查询配置场景
	bdScene := service.DefaultBdSceneService.FindByCh("tianjinmetro")
	if bdScene.ID == 0 {
		return "", errno.ErrNotFound
	}

	// 查询绑定关系或者注册关系
	sceneUser := repository.DefaultBdSceneUserRepository.FindOne(repository.GetSceneUserOne{
		PlatformKey: "tianjinmetro",
		OpenId:      user.OpenId,
	})
	if sceneUser.PlatformUserId == "" {
		app.Logger.Errorf("天津地铁 未找到绑定关系 : %s", user.OpenId)
		return "", errno.ErrBindRecordNotFound
	}

	//请求参数
	Request := MetroRequest{
		AllotId:     "33333333",
		EtUserPhone: sceneUser.Phone,
		AllotNum:    1,
	}

	//加密过程
	jsons, errs := json.Marshal(Request) //转换成JSON返回的是byte[]
	if errs != nil {
		return "", errs
	}
	str := Encode(string(jsons))
	data, _ := encrypt.RsaEncrypt([]byte(str))
	signature := base64.StdEncoding.EncodeToString(data)

	//header头
	options := []httputil.HttpOption{
		httputil.HttpWithHeader("appid", bdScene.AppId),
		httputil.HttpWithHeader("sequence", strconv.FormatInt(time.Now().Unix(), 10)),
		httputil.HttpWithHeader("signature", signature),
	}

	url := bdScene.Domain + "/tj-metro-api/open-forward/api/eTicket/allot"
	body, err := httputil.PostJson(url, Request, options...)
	app.Logger.Infof("天津地铁 返回 : %s", body)
	if err != nil {
		return "", err
	}

	response := MetroResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		app.Logger.Errorf("天津地铁 json_decode_err: %s", err.Error())
		return "", err
	}

	if response.ResultCode != "0000" {
		app.Logger.Errorf("天津地铁 err: %s\n", response.ResultDesc)
		return "", errors.New(response.ResultDesc)
	}
	//记录
	_, err = app.RpcService.CouponRpcSrv.SendCoupon(srv.ctx, &couponclient.SendCouponReq{
		CouponCardTypeId: typeId,
		UserId:           user.ID,
		BizId:            response.ResultData.OrderNo,
		CouponCardTitle:  "天津地铁" + fmt.Sprintf("%.0f", amount) + "元出行红包",
		StartTime:        time.Now().UnixMilli(),
		EndTime:          time.Now().AddDate(0, 0, 90).UnixMilli(),
	})

	if err != nil {
		app.Logger.Errorf("天津地铁 券包 发放错误 : %s\n", err.Error())
		return "", err
	}
	return response.ResultData.OrderNo, nil
}

//参考 https://iswxw.blog.csdn.net/article/details/122612927?spm=1001.2101.3001.6650.4&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7ERate-4-122612927-blog-125201969.pc_relevant_3mothn_strategy_and_data_recovery&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7ERate-4-122612927-blog-125201969.pc_relevant_3mothn_strategy_and_data_recovery&utm_relevant_index=5

func Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
