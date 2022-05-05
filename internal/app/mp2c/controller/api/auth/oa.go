package auth

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model/entity"
	authSrv "mio/internal/pkg/service/auth"
	"mio/internal/pkg/util/apiutil"
	"net/http"
)

var DefaultOaController = OaController{}

type OaController struct {
}

func (OaController) Sign(c *gin.Context) (gin.H, error) {
	form := ConfigSignForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	signResult, err := authSrv.OaService{Platform: entity.UserSource(form.Platform)}.Sign(form.Url)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"sign": signResult,
	}, nil
}

// AutoLogin 手机端完成一次授权后 调用端state未改变 回直接回调到callback而不会经过此方法
func (OaController) AutoLogin(c *gin.Context) {
	form := AutoLoginForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	oaService := authSrv.OaService{
		Platform: entity.UserSource(form.Platform),
	}
	loginUrl, err := oaService.AutoLogin(form.RedirectUri, form.State)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	//使用301重定项时 后续的授权将会已相同的state直接进入callback函数 而不会经过本函数
	//使用302重定向时 后续的授权将会同样先进入本函数 保存数据后进行授权然后进入callback函数
	c.Redirect(http.StatusFound, loginUrl)
	return
}

func (OaController) AutoLoginCallback(c *gin.Context) {
	form := AutoLoginCallbackForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	oaService := authSrv.OaService{}
	redirectUri, err := oaService.AutoLoginCallback(form.Code, form.State)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	//如果使用302重定向 浏览器点击返回上一页时 url会变成此方法url且会重新进入此方法报数据异常错误
	//使用301重定向时 浏览器点击返回上一页时 url会变成此方法url但不会进入此方法 会已相同的code进入回调页
	c.Redirect(http.StatusMovedPermanently, redirectUri)
}

func (OaController) Login(c *gin.Context) (gin.H, error) {
	form := OaAuthForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	oaService := authSrv.OaService{
		Platform: entity.UserSource(form.Platform),
	}
	token, err := oaService.LoginByCode(form.Code)
	return gin.H{
		"token": token,
	}, err
}
