package tianjin_metro

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"math/rand"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/quiz"
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

var channelTypes = map[int64]string{
	1073: "天津地铁",
}

func (srv *Service) SendCoupon(typeId int64, amount float64, user entity.User) (string, error) {
	err := srv.GetTjMetroTicketStatus(user.OpenId, user.ID, user.ChannelId, 1000)
	//调用微服务，发地铁券

	//查询配置场景
	bdScene := service.DefaultBdSceneService.FindByCh("tianjinmetro")
	if bdScene.ID == 0 {
		return "", errno.ErrNotFound
	}

	//请求参数
	Request := MetroRequest{
		AllotId:     "11231231231231",
		EtUserPhone: "15000000000",
		AllotNum:    1,
	}

	//获取签名
	signature, err := getSign(Request)
	if err != nil {
		return "", err
	}

	//header头
	options := []httputil.HttpOption{
		httputil.HttpWithHeader("appid", bdScene.AppId),
		httputil.HttpWithHeader("sequence", getSequence()),
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

func (srv Service) GetTjMetroTicketStatus(openid string, uid int64, channelId int64, typeId int64) error {
	//判断是否指定渠道用户
	_, ok := channelTypes[channelId]
	if !ok {
		return errno.ErrCommon.WithMessage("不满足参与条件")
	}
	//查看是否领取了，没领取满足条件
	couponResp, err := app.RpcService.CouponRpcSrv.FindCoupon(srv.ctx, &couponclient.FindCouponReq{
		UserId:           uid,
		CouponCardTypeId: typeId,
	})
	if err != nil {
		return err
	}
	if couponResp.Exist {
		return errno.ErrCouponReceived
	}
	//查看今天是否答题，没答题满足条件
	availability, err := quiz.DefaultQuizService.Availability(openid)
	if err != nil {
		return err
	}
	if !availability {
		return errno.ErrCommon.WithMessage("不满足答题条件")
	}
	return nil
}

//参考 https://iswxw.blog.csdn.net/article/details/122612927?spm=1001.2101.3001.6650.4&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7ERate-4-122612927-blog-125201969.pc_relevant_3mothn_strategy_and_data_recovery&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7ERate-4-122612927-blog-125201969.pc_relevant_3mothn_strategy_and_data_recovery&utm_relevant_index=5

/*func Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}*/

//签名开始
/**
转换：将请求参数转换为json消息。
摘要：把转换好的字符串采用utf-8编码，使用摘要算法对编码后的字节流进行摘要。使用MD5算法，对转换后的字符串进行摘要，如：md5(json)；
将摘要得到的字节流结果使用十六进制大写表示，如：hex(“helloworld”.getBytes(“utf-8”)) = “68656C6C6F776F726C64”
签名：使用加密算法对摘要后的16进制文本进行加密。
使用RSA算法，1024位，填充方式采用RSA/ECB/PKCS1Padding，如RSA(“ 68656C6C6F776F726C64”, key)。
*/

func getSign(request MetroRequest) (string, error) {
	jsons, errs := json.Marshal(request) //转换成JSON返回的是byte[]
	if errs != nil {
		return "", errs
	}

	h := md5.New()
	h.Write(jsons)
	str := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

	data, _ := encrypt.RsaEncrypt([]byte(str))
	signature := base64.StdEncoding.EncodeToString(data) //签名
	return signature, nil
}

//业务流水  yyyyMMddHHmmss+10位数字（在一定时间内不重复），仅作为接口调用跟踪用途，不作为业务用途，业务流水在业务接口中定义。

func getSequence() string {
	timeNowStr := time.Now().Format("20060102150405")
	rand.Seed(time.Now().Unix())                                          //Seed生成的随机数
	sequence := timeNowStr + strconv.Itoa(random(1000000000, 9999999999)) //业务流水
	return sequence
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
