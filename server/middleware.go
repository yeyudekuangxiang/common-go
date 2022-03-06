package server

import (
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"log"
	"mio/core/app"
	"mio/internal/errno"
	"mio/internal/util"
	"mio/internal/zap"
	"mio/model/entity"
	"mio/service"
	"time"
)

func Middleware(middleware *gin.Engine) {
	middleware.Use(corsM())
	middleware.Use(gin.Recovery())
	middleware.Use(access())
}
func access() gin.HandlerFunc {
	//执行测试时 访问日志输出到控制台
	if util.IsTesting() {
		return ginzap.Ginzap(zap.DefaultLogger("info"), time.RFC3339, false)
	}
	logger := zap.NewZapLogger(zap.LoggerConfig{
		Level:    "info",
		Path:     "runtime",
		FileName: "access.log",
		MaxSize:  100,
	})
	return ginzap.Ginzap(logger, time.RFC3339, false)
}
func auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			//ctx.AbortWithStatusJSON(200,http.ErrorResponse(http.NewBaseError(400,"未登录")))
			ctx.Set("AuthUser", entity.User{})
			return
		}

		user, err := service.DefaultUserService.GetUserByToken(token)
		if err != nil || user == nil {
			//ctx.AbortWithStatusJSON(200,http.ErrorResponse(http.NewBaseError(400,err.Error())))
			ctx.Set("AuthUser", entity.User{})
			return
		}
		ctx.Set("AuthUser", *user)
	}
}
func authAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			ctx.AbortWithStatusJSON(200, formatErr(errno.ErrAuth, nil))
			return
		}

		admin, err := service.DefaultAdminService.GetAdminByToken(token)
		if err != nil || admin == nil {
			app.Logger.Error("用户登陆验证失败", admin, err)
			ctx.AbortWithStatusJSON(200, formatErr(errno.ErrValidation, nil))
			return
		}

		ctx.Set("AuthAdmin", *admin)
	}
}
func mustAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			ctx.AbortWithStatusJSON(200, formatErr(errno.ErrAuth, nil))
			return
		}

		user, err := service.DefaultUserService.GetUserByToken(token)
		if err != nil || user == nil {
			app.Logger.Error("用户登陆验证失败", user, err)
			ctx.AbortWithStatusJSON(200, formatErr(errno.ErrValidation, nil))
			return
		}
		ctx.Set("AuthUser", *user)
	}
}

//临时使用openid作为登陆验证
func mustAuth2() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		openId := ctx.GetHeader("openid")
		if openId == "" {
			ctx.AbortWithStatusJSON(200, formatErr(errno.ErrAuth, nil))
			return
		}

		user, err := service.DefaultUserService.GetUserByOpenId(openId)
		if err != nil || user == nil {
			app.Logger.Error("用户登陆验证失败", user, err)
			ctx.AbortWithStatusJSON(200, formatErr(errno.ErrValidation, nil))
			return
		}
		ctx.Set("AuthUser", *user)
	}
}

//临时使用openid作为登陆验证
func auth2() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		openId := ctx.GetHeader("openid")
		if openId == "" {
			ctx.Set("AuthUser", entity.User{})
			return
		}

		user, err := service.DefaultUserService.GetUserByOpenId(openId)
		if err != nil || user == nil {
			ctx.Set("AuthUser", entity.User{})
			return
		}
		ctx.Set("AuthUser", *user)
	}
}

type ThrottleConfig struct {
	Throttle string
}

func throttle() gin.HandlerFunc {
	throttleConfig := &struct {
		Throttle string
	}{}
	_ = app.Ini.Section("http").MapTo(throttleConfig)
	if throttleConfig.Throttle == "" {
		throttleConfig.Throttle = "200-M"
	}
	rate, err := limiter.NewRateFromFormatted(throttleConfig.Throttle)
	if err != nil {
		log.Fatal(err)
	}

	store := memory.NewStoreWithOptions(limiter.StoreOptions{
		Prefix: "throttle",
	})

	middleware := mgin.NewMiddleware(limiter.New(store, rate), mgin.WithKeyGetter(func(c *gin.Context) string {
		return util.Md5(c.ClientIP() + c.Request.Method + c.FullPath())
	}))
	return middleware
}
func corsM() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("x-token", "token", "authorization", "openid")
	return cors.New(config)
}
