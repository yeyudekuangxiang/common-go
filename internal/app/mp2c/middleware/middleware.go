package middleware

import (
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"log"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	service2 "mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"mio/pkg/zap"
	"time"
)

// Middleware 全局中间件
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

		user, err := service2.DefaultUserService.GetUserByToken(token)
		if err != nil {
			//ctx.AbortWithStatusJSON(200,http.ErrorResponse(http.NewBaseError(400,err.Error())))
			ctx.Set("AuthUser", entity.User{})
			return
		}
		ctx.Set("AuthUser", *user)
	}
}
func AuthAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			ctx.AbortWithStatusJSON(200, util.FormatErr(errno.ErrAuth, nil))
			return
		}

		admin, err := service2.DefaultAdminService.GetAdminByToken(token)
		if err != nil || admin == nil {
			app.Logger.Error("用户登陆验证失败", admin, err)
			ctx.AbortWithStatusJSON(200, util.FormatErr(errno.ErrValidation, nil))
			return
		}

		ctx.Set("AuthAdmin", *admin)
	}
}
func mustAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			ctx.AbortWithStatusJSON(200, util.FormatErr(errno.ErrAuth, nil))
			return
		}

		user, err := service2.DefaultUserService.GetUserByToken(token)
		if err != nil || user.ID == 0 {
			app.Logger.Error("用户登陆验证失败", user, err)
			ctx.AbortWithStatusJSON(200, util.FormatErr(errno.ErrValidation, nil))
			return
		}
		ctx.Set("AuthUser", *user)
	}
}

//临时使用openid作为登陆验证
func MustAuth2() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user *entity.User
		var err error
		if token := ctx.GetHeader("token"); token != "" {
			user, err = service2.DefaultUserService.GetUserByToken(token)
			if err != nil || user.ID == 0 {
				app.Logger.Error("mustAuth token err", token, err)
				ctx.AbortWithStatusJSON(200, util.FormatErr(errno.ErrValidation, nil))
				return
			}
		}

		if openId := ctx.GetHeader("openid"); openId != "" {
			user, err = service2.DefaultUserService.GetUserByOpenId(openId)
			if err != nil || user == nil {
				app.Logger.Error("mustAuth openid err", openId, err)
				ctx.AbortWithStatusJSON(200, util.FormatErr(errno.ErrValidation, nil))
				return
			}
		}

		if user == nil {
			ctx.AbortWithStatusJSON(200, util.FormatErr(errno.ErrAuth, nil))
			return
		}
		ctx.Set("AuthUser", *user)
	}
}

//临时使用openid作为登陆验证
func Auth2() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user *entity.User
		var err error

		if token := ctx.GetHeader("token"); token != "" {
			user, err = service2.DefaultUserService.GetUserByToken(token)
			if err != nil {
				app.Logger.Error("auth token err", token, err)
			}
		}

		if openId := ctx.GetHeader("openid"); openId != "" {
			user, err = service2.DefaultUserService.GetUserByOpenId(openId)
			if err != nil {
				app.Logger.Error("auth openid err", openId, err)
			}
		}
		if user == nil {
			user = &entity.User{}
		}
		ctx.Set("AuthUser", *user)
	}
}

type ThrottleConfig struct {
	Throttle string
}

func Throttle() gin.HandlerFunc {
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