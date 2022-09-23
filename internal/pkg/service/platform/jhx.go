package platform

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NewJhxService(ctx *context.MioContext) *JhxService {
	return &JhxService{
		ctx:    ctx,
		Domain: "http://m.jinhuaxing.com.cn/api",
		CommonRequest: commonRequest{
			AppId:     "2498728d209d",
			Version:   "1.0",
			Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
			Nonce:     strconv.Itoa(rand.Int()),
		},
	}
}

type JhxService struct {
	ctx           *context.MioContext
	Domain        string `json:"domain"`
	CommonRequest commonRequest
}

type commonRequest struct {
	AppId     string `json:"appid"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
}

func (srv JhxService) TicketCreate(tradeno string, user entity.User) error {
	params := make(map[string]interface{}, 0)
	_ = util.MapTo(srv.CommonRequest, &params)
	params["tradeno"] = tradeno
	params["mobile"] = user.PhoneNumber
	sign := srv.getSign(params)
	params["sign"] = strings.ToUpper(sign)
	url := srv.Domain + "/busticket/ticket_create"
	body, err := httputil.PostJson(url, params)
	fmt.Printf("%s\n", body)
	if err != nil {
		return err
	}
	response := jhxCommonResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
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
	//入库
	//ticketCreateResponse.QrCodeStr
	fmt.Printf("%v\n", response)
	return nil
}

//消费通知
func (srv JhxService) TicketNotify(sign string, params map[string]interface{}) error {
	md5Sign := srv.getSign(params)
	if sign != md5Sign {
		return errors.New("验签失败")
	}
	//查询库 根据tradeno获取券码

	//如果 status 相等 不处理 返回 nil

	//如果 status 不想等 根据 tadeno 更新status,used_time 返回nil

	//如果有err 返回err
	return nil
}

func (srv JhxService) TicketStatus(tradeno string) (*jhxTicketStatusResponse, error) {
	params := make(map[string]string, 0)
	params["tradeno"] = tradeno
	url := srv.Domain + "/busticket/ticket_create"
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

// GetSign 签名
func (srv JhxService) getSign(params map[string]interface{}) string {
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += v + "=" + util.InterfaceToString(params[v]) + "&"
	}
	signStr = strings.TrimRight(signStr, "&")
	return encrypt.Md5(signStr)
}
