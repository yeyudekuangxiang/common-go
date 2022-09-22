package platform

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"mio/internal/pkg/core/context"
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

type senCouponResponse struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Time string                 `json:"time"`
	Data map[string]interface{} `json:"data"`
}

func (srv JhxService) TicketCreate(tradeno, mobile string) error {
	params := make(map[string]interface{}, 0)
	_ = util.MapTo(srv.CommonRequest, &params)
	params["tradeno"] = tradeno
	params["mobile"] = mobile
	sign := srv.getSign(params)
	params["sign"] = strings.ToUpper(sign)
	url := srv.Domain + "/busticket/ticket_create"
	body, err := httputil.PostJson(url, params)
	fmt.Printf("%s\n", body)
	if err != nil {
		return err
	}
	response := senCouponResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Code != 0 {
		return errors.New(response.Msg)
	}
	fmt.Printf("%v\n", response)
	return nil
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

//消费通知
func (srv JhxService) TicketNotify() {

}
