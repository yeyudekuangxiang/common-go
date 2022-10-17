package ytx

import (
	"encoding/json"
	"errors"
	"fmt"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	platformUtil "mio/pkg/platform"
	"sort"
	"strings"
	"time"
)

type ytxOption struct {
	Domain string
	Secret string
}

type Options func(options *ytxOption)

//openid:  CpziorTGUL02NrrBqsbbhsAN0Ve4ZMSpPEmgBPAGZOY=
//secret:   test_secret

func NewYtxService(ctx *context.MioContext, jhxOptions ...Options) *Service {
	options := &ytxOption{
		Domain: " https://apift.ruubypay.com",
	}

	for i := range jhxOptions {
		jhxOptions[i](options)
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

func (srv *Service) Synchro(memberId string, openId string) error {
	synchroRequest := SynchroRequest{
		OpenId:         memberId,
		RegDate:        time.Now().Format("20060102150405"),
		PlatformUserId: openId,
		Ts:             time.Now().UnixMilli(),
	}
	params := make(map[string]interface{}, 0)
	err := util.MapTo(&synchroRequest, &params)
	if err != nil {
		return err
	}
	params["secret"] = srv.option.Secret
	synchroRequest.Signature = platformUtil.GetSign(params, "", "&")
	url := srv.option.Domain + "/markting_activity/network/lvmiao/synchro"
	body, err := httputil.PostJson(url, synchroRequest)
	fmt.Printf("ytx synchro response body: %s\n", body)
	if err != nil {
		return err
	}
	response := synchroResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Unmarshal body: %s\n", err.Error())
		return err
	}
	if response.ResCode != "0000" {
		return errors.New(response.ResMessage)
	}
	return nil
}

func (srv *Service) getSign(synchroRequest SynchroRequest) string {
	params := make(map[string]interface{}, 0)
	err := util.MapTo(&synchroRequest, &params)
	if err != nil {
		return ""
	}
	params["secret"] = srv.option.Secret

	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += v + "=$" + util.InterfaceToString(params[v]) + "&"
	}
	signStr = strings.TrimRight(signStr, "&")
	//验证签名
	return encrypt.Md5(signStr)
}
