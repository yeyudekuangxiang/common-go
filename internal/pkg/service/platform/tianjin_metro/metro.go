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
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	"mio/pkg/errno"
	"strconv"
	"strings"
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

	grantV5Request := GrantV4Request{
		AllotId:     "33333333",
		EtUserPhone: sceneUser.Phone,
		AllotNum:    1,
	}
	//加密过程
	jsons, errs := json.Marshal(grantV5Request) //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println(string(jsons))

	str := Encode(string(jsons))
	data, _ := encrypt.RsaEncrypt([]byte(str))
	signature := base64.StdEncoding.EncodeToString(data)

	authToken := httputil.HttpWithHeader("appid", "264735a59163453d9772f92e1f703123")         //天津地铁分配给开发者/商户的appid
	authToken = httputil.HttpWithHeader("random", "111")                                      //加密后的随机码，当报文中有需要加密的字段的时候需要传此参数
	authToken = httputil.HttpWithHeader("sequence", strconv.FormatInt(time.Now().Unix(), 10)) //yyyyMMddHHmmss+10位数字（在一定时间内不重复），仅作为接口调用跟踪用途，不作为业务用途，业务流水在业务接口中定义。
	authToken = httputil.HttpWithHeader("version", "1.0")                                     //版本号 1.0。
	authToken = httputil.HttpWithHeader("signature", signature)                               //签名

	url := "https://app.trtpazyz.com/tj-metro-api/open-forward/api/eTicket/allot"
	body, err := httputil.PostJson(url, grantV5Request, authToken)
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

	return response.SubData.OrderNo, nil

	//return "", nil
}

func Check(content, encrypted string) bool {
	return strings.EqualFold(Encode(content), encrypted)
}

//学习网址 https://iswxw.blog.csdn.net/article/details/122612927?spm=1001.2101.3001.6650.4&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7ERate-4-122612927-blog-125201969.pc_relevant_3mothn_strategy_and_data_recovery&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7ERate-4-122612927-blog-125201969.pc_relevant_3mothn_strategy_and_data_recovery&utm_relevant_index=5

func Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// RSA算法

/*
func (srv *Service) getAppSecret() string {
	t := time.Now().Unix()
	return encrypt.Md5(srv.option.AppId + encrypt.Md5(srv.option.Secret) + strconv.FormatInt(t, 10))
}*/
