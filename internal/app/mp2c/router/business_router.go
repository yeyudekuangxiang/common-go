package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/business"
	"mio/internal/app/mp2c/middleware"
	"mio/internal/pkg/util/apiutil"
)

func BusinessRouter(router *gin.Engine) {
	businessRouter := router.Group("/api/mp2c/business")
	businessRouter.Use(middleware.Throttle())
	{

	}

	authRouter := router.Group("/api/mp2c/business")
	authRouter.Use(middleware.Throttle())
	{
		carbonRouter := authRouter.Group("/carbon")
		{
			carbonRouter.GET("/record/list", apiutil.Format(business.DefaultCarbonCreditsController.GetCarbonCreditLogInfoList))

			carbonRouter.GET("/rank/user/list", apiutil.Format(business.DefaultCarbonRankController.GetUserRankList))
		}

		pointRouter := authRouter.Group("/point")
		{
			pointRouter.GET("/record/list", apiutil.Format(business.DefaultPointController.GetPointRecordList))
		}
	}
}
