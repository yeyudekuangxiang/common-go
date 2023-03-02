package {{.pkgName}}

import (
	{{.imports}}
	"net/http"
)

type {{.logic}} struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	req *http.Request
}

func New{{.logic}}(ctx context.Context, svcCtx *svc.ServiceContext,req *http.Request) *{{.logic}} {
	return &{{.logic}}{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		req: req,
	}
}

func (l *{{.logic}}) {{.function}}({{.request}}) {{.responseType}} {
	// todo: add your logic here and delete this line

	{{.returnString}}
}
