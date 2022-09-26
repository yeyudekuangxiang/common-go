package platform

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math/rand"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
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

func (srv JhxService) TicketCreate(tradeno string, user entity.User) error {
	params := srv.getCommonParams()
	params["tradeno"] = tradeno
	params["mobile"] = user.PhoneNumber
	sign := srv.getJhxSign(params)
	params["sign"] = strings.ToUpper(sign)
	url := srv.option.Domain + "/busticket/ticket_create"
	body, err := httputil.PostJson(url, params)
	fmt.Printf("response body: %s\n", body)
	if err != nil {
		return err
	}
	//response := jhxCommonResponse{}
	response := make(map[string]interface{}, 0)
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Unmarshal body: %s\n", err.Error())
		return err
	}
	//if response.Code != 0 {
	//	return errors.New(response.Msg)
	//}
	//ticketCreateResponse := &jhxTicketCreateResponse{}
	//err = util.MapTo(response.Data, &ticketCreateResponse)
	//if err != nil {
	//	return err
	//}
	// todo 入库
	//ticketCreateResponse.QrCodeStr
	fmt.Printf("%v\n", response)
	return nil
}

//消费通知
func (srv JhxService) TicketNotify(sign string, params map[string]string) error {
	if err := srv.checkSign(sign, params); err != nil {
		return err
	}
	//查询库 根据tradeno获取券码

	//如果 status 相等 不处理 返回 nil

	//如果 status 不想等 根据 tadeno 更新status,used_time 返回nil

	//如果有err 返回err
	return nil
}

func (srv JhxService) TicketStatus(tradeno string) (*jhxTicketStatusResponse, error) {
	params := srv.getCommonParams()
	params["tradeno"] = tradeno

	sign := srv.getJhxSign(params)
	params["sign"] = strings.ToUpper(sign)

	url := srv.option.Domain + "/busticket/ticket_create"
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

func (srv JhxService) PreCollectPoint(sign string, params map[string]string) error {
	if err := srv.checkSign(sign, params); err != nil {
		return err
	}
	//根据 platform_member_id 获取 openid
	sceneUser := repository.DefaultBdSceneUserRepository.FindPlatformUserByPlatformUserId(params["memberId"], "jhx")
	if sceneUser.ID == 0 {
		return errors.New("未找到绑定关系")
	}
	//创建数据
	fromString, err := decimal.NewFromString(params["amount"])
	if err != nil {
		return err
	}
	point := fromString.Mul(decimal.NewFromInt(10)).Round(2).String()
	err = repository.DefaultBdScenePrePointRepository.Create(entity.BdScenePrePoint{
		PlatformKey:    sceneUser.PlatformKey,
		PlatformUserId: sceneUser.PlatformUserId,
		Point:          point,
		OpenId:         sceneUser.OpenId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (srv JhxService) GetPreCollectPointList(sign string, params map[string]string) ([]entity.BdScenePrePoint, error) {
	if err := srv.checkSign(sign, params); err != nil {
		return nil, err
	}
	//根据 platform_member_id 获取 openid
	sceneUser := repository.DefaultBdSceneUserRepository.FindPlatformUserByPlatformUserId(params["memberId"], "jhx")
	if sceneUser.ID == 0 {
		return nil, errors.New("未找到绑定关系")
	}
	//获取pre_point数据
	list := repository.DefaultBdScenePrePointRepository.FindBy(repository.GetScenePrePoint{
		PlatformKey:    sceneUser.PlatformKey,
		PlatformUserId: sceneUser.PlatformUserId,
		OpenId:         sceneUser.OpenId,
		StartTime:      time.Now().AddDate(0, 0, -7),
		EndTime:        time.Now(),
	})
	return list, nil
}

//消费数据
func (srv JhxService) CollectPoint(sign string, params map[string]string) error {
	if err := srv.checkJhxSign(sign, params); err != nil {
		return err
	}
	//根据 platform_member_id 获取 openid

	//获取pre_point数据 one limit

	//调用point_trans.incPoint

	//删除pre_point对应数据

	return nil
}

func (srv JhxService) checkJhxSign(sign string, params map[string]string) error {
	md5Sign := srv.getJhxSign(params)
	if sign != md5Sign {
		return errors.New("验签失败")
	}
	return nil
}

func (srv JhxService) checkSign(sign string, params map[string]string) error {
	md5Sign := srv.getSign(params)
	if sign != md5Sign {
		return errors.New("验签失败")
	}
	return nil
}

// GetSign 签名
func (srv JhxService) getJhxSign(params map[string]string) string {
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
