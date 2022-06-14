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
	authRouter.Use(middleware.Throttle(), middleware.AuthBusinessUser())
	{
		carbonRouter := authRouter.Group("/carbon")
		{
			carbonRouter.GET("/record/list", apiutil.Format(business.DefaultCarbonCreditsController.GetCarbonCreditLogInfoList))
			carbonRouter.GET("/rank/user/list", apiutil.Format(business.DefaultCarbonRankController.GetUserRankList))
			carbonRouter.GET("/rank/department/list", apiutil.Format(business.DefaultCarbonRankController.GetDepartmentRankList))
			carbonRouter.POST("/rank/user/like/status/change", apiutil.Format(business.DefaultCarbonRankController.ChangeUserRankLikeStatus))
			carbonRouter.POST("/rank/department/like/status/change", apiutil.Format(business.DefaultCarbonRankController.ChangeDepartmentRankLikeStatus))

			carbonRouter.POST("/collect/evcar", apiutil.Format(business.DefaultCarbonController.CollectEvCar))
			carbonRouter.POST("/collect/online-meeting", apiutil.Format(business.DefaultCarbonController.CollectOnlineMeeting))
			carbonRouter.POST("/collect/public-transport", apiutil.Format(business.DefaultCarbonController.CollectPublicTransport))
			carbonRouter.POST("/collect/save-water-electricity", apiutil.Format(business.DefaultCarbonController.CollectSaveWaterElectricity))
			//获取企业的本月减碳来源列表（按本月减碳量排序）（任务3
			carbonRouter.GET("/sorted/list", apiutil.Format(business.DefaultCarbonCreditsController.GetCarbonCreditLogSortedList))
			//9 返回当前用户减碳来源以及本月减碳数量（按减碳量排序）（任务3）
			carbonRouter.GET("/user/sorted/list", apiutil.Format(business.DefaultCarbonCreditsController.GetUserCarbonCreditLogSortedList))
			//10 获取企业近6个月每个月减碳数量列表（任务3）
			//11 获取企业减碳场景以及每个场景最近六个月减碳数量总和 最多20个减碳场景（任务3）
			carbonRouter.GET("/record/history", apiutil.Format(business.DefaultCarbonCreditsController.GetCarbonCreditLogSortedListHistory))
			//用户历史上每个月的减碳情况
			carbonRouter.GET("/record/user/history", apiutil.Format(business.DefaultCarbonCreditsController.GetCarbonCreditLogSortedListUserHistory))
		}

		pointRouter := authRouter.Group("/point")
		{
			pointRouter.GET("/record/list", apiutil.Format(business.DefaultPointController.GetPointRecordList))
		}
		userRouter := authRouter.Group("/user")
		{
			//7 获取当前用户的信息（任务1）
			//8 获取当前用户本月减碳信息(任务3)
			userRouter.GET("/info", apiutil.Format(business.DefaultUserController.GetUserInfo))
		}

		CompanyRouter := authRouter.Group("/company")
		{
			//获取企业信息以及碳减排信息（任务1）
			CompanyRouter.GET("/info", apiutil.Format(business.DefaultCompanyController.GetCompanyInfo))
		}
	}
}
