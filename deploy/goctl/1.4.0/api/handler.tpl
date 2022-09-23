package {{.PkgName}}

import (
	"net/http"
	"{{.projectPath}}/common/result"
	"{{.projectPath}}/common/tool/api"

	"github.com/zeromicro/go-zero/rest/httpx"
	{{.ImportPackages}}
)

func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := api.BindForm(r, &req); err != nil {
			result.HttpResult(r, w, nil, err)
			return
		}

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		{{if .HasResp}}result.HttpResult(r, w, resp, err){{else}}result.HttpResult(r, w, nil, err){{end}}
	}
}
