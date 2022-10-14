package ytx

import (
	"mio/internal/pkg/core/context"
)

type jhxOption struct {
	Domain string
	Path   string
	Secret string
}

type JhxOptions func(options *jhxOption)

func NewJhxService(ctx *context.MioContext, jhxOptions ...JhxOptions) *JhxService {
	options := &jhxOption{
		Domain: "http://m.jinhuaxing.com.cn/api",
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
