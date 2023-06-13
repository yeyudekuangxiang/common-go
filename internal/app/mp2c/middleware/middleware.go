package middleware

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	mzap "gitlab.miotech.com/miotech-application/backend/common-go/logger/zap"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/wxwork"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/queue/producer/wxworkpdr"
	"mio/internal/pkg/queue/types/message/wxworkmsg"
	service2 "mio/internal/pkg/service"
	"mio/internal/pkg/service/business"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"os"
	"runtime"
	"strings"
	"time"
)

// Middleware 全局中间件
func Middleware(middleware *gin.Engine) {
	middleware.Use(corsM())
	middleware.Use(recovery())
	middleware.Use(access())
}
func recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		e, ok := err.(error)
		if ok {
			c.JSON(apiutil.FormatErr(e, nil))
		} else {
			c.JSON(200, apiutil.FormatResponse(errno.ErrInternalServer.Code(), nil, fmt.Sprintf("%v", err)))
		}
		c.Abort()

		callers := callers()
		go func() {
			if config.Config.App.Env != "prod" {
				return
			}

			sendErr := wxworkpdr.SendRobotMessage(wxworkmsg.RobotMessage{
				Key:  config.Constants.WxWorkBugRobotKey,
				Type: wxwork.MsgTypeMarkdown,
				Message: wxwork.Markdown{
					Content: fmt.Sprintf("**容器:**%s \n\n**来源:**panic \n\n**消息:**%+v \n\n**堆栈:**%s \n\n<@all>", os.Getenv("HOSTNAME"), err, callers),
				},
			})

			if sendErr != nil {
				log.Printf("推送异常到企业微信失败 %v %v", err, sendErr)
			}
		}()
	})
}
func callers() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(5, pcs[:])

	s := strings.Builder{}
	for _, pc := range pcs[0:n] {
		f := errors.Frame(pc)
		s.WriteString(fmt.Sprintf("\n%+v", f))
	}
	return s.String()
}

func access() gin.HandlerFunc {
	return Access(app.OriginLogger.With(mzap.LogAccess), time.RFC3339, false)
}

func AuthAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			ctx.AbortWithStatusJSON(apiutil.FormatErr(errno.ErrAuth, nil))
			return
		}

		admin, exists, err := service2.DefaultSystemAdminService.GetAdminByToken(token)
		if err != nil || !exists {
			app.Logger.Error("用户登陆验证失败", admin, err)
			ctx.AbortWithStatusJSON(apiutil.FormatErr(err, nil))
			return
		}
		ctx.Set("AuthAdmin", *admin)
	}
}

// AuthBusinessUser 企业用户登录时,解析token
func AuthBusinessUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("b-token")
		if token == "" {
			ctx.AbortWithStatusJSON(apiutil.FormatErr(errno.ErrAuth, nil))
			return
		}

		user, err := business.DefaultUserService.GetBusinessUserByToken(token)
		if err != nil || user == nil {
			app.Logger.Error("用户登陆验证失败", user, err)
			ctx.AbortWithStatusJSON(apiutil.FormatErr(err, nil))
			return
		}

		ctx.Set("BusinessUser", *user)
	}
}

func MustAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user *entity.User
		var err error
		if token := ctx.GetHeader("token"); token != "" {
			user, err = service2.DefaultUserService.GetUserByToken(token)
			if err != nil || user.ID == 0 {
				app.Logger.Error("mustAuth token err", token, err)
				ctx.AbortWithStatusJSON(apiutil.FormatErr(errno.ErrValidation.WithErrMessage(token), nil))
				return
			}
		}

		if user == nil {
			ctx.AbortWithStatusJSON(apiutil.FormatErr(errno.ErrAuth.WithCaller(), nil))
			return
		}
		if user.UserStatus != 1 {
			ctx.AbortWithStatusJSON(apiutil.FormatErr(errno.ErrAuth.WithCaller(), nil))
			return
		}

		ctx.Set("AuthUser", *user)
	}
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user *entity.User
		var err error

		if token := ctx.GetHeader("token"); token != "" {
			user, err = service2.DefaultUserService.GetUserByToken(token)
			if err != nil {
				app.Logger.Error("auth token err", token, err)
			}
		}

		if user == nil {
			user = &entity.User{}
		} else if user.UserStatus != 1 {
			ctx.AbortWithStatusJSON(apiutil.FormatErr(errno.ErrAuth.WithCaller(), nil))
			return
		}

		ctx.Set("AuthUser", *user)
	}
}

func Throttle() gin.HandlerFunc {

	throttle := "200-M"
	if config.Config.Http.Throttle != "" {
		throttle = config.Config.Http.Throttle
	}

	rate, err := limiter.NewRateFromFormatted(throttle)
	if err != nil {
		log.Fatal(err)
	}

	store, err := redis.NewStoreWithOptions(app.Redis, limiter.StoreOptions{
		Prefix: "mp2c:throttle",
	})
	if err != nil {
		log.Fatal("创建limit失败", err)
	}

	middleware := mgin.NewMiddleware(limiter.New(store, rate), mgin.WithKeyGetter(func(c *gin.Context) string {
		return encrypttool.Md5(c.ClientIP() + c.Request.Method + c.FullPath())
	}))
	return middleware
}
func corsM() gin.HandlerFunc {
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AddAllowHeaders("x-token", "b-token", "token", "authorization", "openid", "channel")
	return cors.New(cfg)
}

// 临时使用openid作为登陆验证
func MqAuth2() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			ctx.AbortWithStatusJSON(apiutil.FormatErr(errno.ErrAuth, nil))
			return
		}
		if token != config.FindMqToken(ctx.FullPath()) {
			ctx.AbortWithStatusJSON(apiutil.FormatErr(errno.ErrValidation.WithErrMessage(token), nil))
			return
		}
	}
}
