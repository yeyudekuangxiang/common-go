package ccring

import (
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"mio/internal/pkg/core/app"
)

type ccRingService struct {
	Authorization string `json:"authorization"`
	Url           string `json:"interface"`
	Domain        string `json:"domain"`
	Option        *ccRingOption
}

type ccRingOption struct {
	MemberId            string  `json:"memberId"`
	OrderNum            string  `json:"orderNum"`
	DegreeOfCharge      float64 `json:"degreeOfCharge,omitempty"`
	ProductCategoryName string  `json:"productCategoryName,omitempty"`
	Name                string  `json:"name,omitempty"`
	Qua                 string  `json:"qua,omitempty"`
}

type CcRingOptions func(option *ccRingOption)

func NewCCRingService(token, domain, url string, opts ...CcRingOptions) *ccRingService {
	options := &ccRingOption{}
	for i := range opts {
		opts[i](options)
	}
	return &ccRingService{
		Authorization: token,
		Domain:        domain,
		Url:           url,
		Option:        options,
	}
}

func WithCCRingMemberId(memberId string) CcRingOptions {
	return func(option *ccRingOption) {
		option.MemberId = memberId
	}
}
func WithCCRingDegreeOfCharge(degree float64) CcRingOptions {
	return func(option *ccRingOption) {
		option.DegreeOfCharge = degree
	}
}

func WithCCRingOrderNum(orderNum string) CcRingOptions {
	return func(option *ccRingOption) {
		option.OrderNum = orderNum
	}
}

func WithCCRingProductCategoryName(categoryName string) CcRingOptions {
	return func(option *ccRingOption) {
		option.ProductCategoryName = categoryName
	}
}
func WithCCRingName(name string) CcRingOptions {
	return func(option *ccRingOption) {
		option.Name = name
	}
}
func WithCCRingQua(qua string) CcRingOptions {
	return func(option *ccRingOption) {
		option.Qua = qua
	}
}

//回调ccring
func (srv ccRingService) CallBack() (string, error) {
	//回调
	authToken := httptool.HttpWithHeader("Authorization", srv.Authorization)
	body, err := httptool.PostJson(srv.Domain+srv.Url, srv.Option, authToken)
	if err != nil {
		app.Logger.Errorf("回调光环错误 url:%v option: %+v Authorization: %v body: %v error: %+v", srv.Domain+srv.Url, srv.Option, srv.Authorization, string(body), err)
		return fmt.Sprintf("%s", body), err
	}
	return fmt.Sprintf("%s", body), nil
}
