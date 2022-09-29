package ccring

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util/httputil"
)

type ccRingOption struct {
	MemberId            string  `json:"memberId"`
	DegreeOfCharge      float64 `json:"degreeOfCharge,omitempty"`
	ProductCategoryName string  `json:"productCategoryName,omitempty"`
	Name                string  `json:"name,omitempty"`
	Qua                 string  `json:"qua,omitempty"`
}

type ccRingService struct {
	Authorization string `json:"authorization"`
	Interface     string `json:"interface"`
	option        *ccRingOption
}

type CcRingOptions func(option *ccRingOption)

func NewCCRingService(token, face string, opts ...CcRingOptions) *ccRingService {
	options := &ccRingOption{}
	for i := range opts {
		opts[i](options)
	}
	return &ccRingService{
		Authorization: token,
		Interface:     face,
		option:        options,
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
func (srv ccRingService) CallBack() {
	scene := repository.DefaultBdSceneRepository.FindByCh("ccring")
	if scene.ID == 0 {
		app.Logger.Errorf("回调光环错误:%s", "未设置scene")
		return
	}
	authToken := httputil.HttpWithHeader("Authorization", srv.Authorization)
	_, err := httputil.PostJson(srv.Interface, srv.option, authToken)
	if err != nil {
		app.Logger.Errorf("回调光环错误:%s", err.Error())
		return
	}
	return
}
