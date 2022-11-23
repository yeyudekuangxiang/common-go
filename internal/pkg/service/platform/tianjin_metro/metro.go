package tianjin_metro

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/coupon"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"math/rand"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
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

func (srv *Service) SendCoupon(typeId int64, user entity.User) (*coupon.SendCouponV2Resp, error) {
	_, err := srv.GetTjMetroTicketStatus(config.Config.ThirdCouponTypes.TjMetro, user.OpenId)
	if err != nil {
		return nil, err
	}
	//调用微服务，发地铁券
	SendCouponV2Resp, err := app.RpcService.CouponRpcSrv.SendCouponV2(srv.ctx, &couponclient.SendCouponV2Req{
		UserId:              user.ID,
		ThirdUserId:         user.PhoneNumber,
		CouponCardTypeId:    typeId,
		BizId:               util.UUID(),
		DistributionChannel: "天津地铁答题发电子票",
	})
	if err != nil {
		return nil, err
	}
	return SendCouponV2Resp, nil
}

func (srv Service) GetTjMetroTicketStatus(typeId int64, openid string) (*entity.User, error) {
	userInfo, exit, _ := repository.DefaultUserRepository.GetUser(repository.GetUserBy{
		OpenId: openid,
	})
	//判断是否注册绿喵
	if !exit {
		app.Logger.Errorf("天津地铁 未注册到绿喵平台 : %s", openid)
		return nil, errno.ErrBindRecordNotFound
	}
	//判断是否指定渠道用户
	_, ok := channelTypes[userInfo.ChannelId]
	if !ok {
		return nil, errno.ErrChannelErr
	}
	//查看是否领取了，没领取满足条件
	couponResp, err := app.RpcService.CouponRpcSrv.FindCoupon(srv.ctx, &couponclient.FindCouponReq{
		UserId:           userInfo.ID,
		CouponCardTypeId: typeId,
	})
	if err != nil {
		return nil, err
	}
	if couponResp.Exist {
		return nil, errno.ErrCouponReceived
	}
	//查看今天是否答题，没答题满足条件
	/*availability, err := quiz.DefaultQuizService.Availability(openid)
	if err != nil {
		return nil, err
	}
	if !availability {
		return nil, errno.ErrCommon.WithMessage("不满足答题条件")
	}*/
	return userInfo, nil
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
