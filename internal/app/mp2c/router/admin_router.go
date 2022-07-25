package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/admin"
	"mio/internal/app/mp2c/middleware"
	"mio/internal/pkg/util/apiutil"
)

func adminRouter(router *gin.Engine) {

	router.POST("/admin/login", apiutil.Format(admin.DefaultAdminController.Login))
	adminRouter := router.Group("/admin")
	adminRouter.Use(middleware.AuthAdmin())
	{

		adminRouter.GET("/info/list", apiutil.Format(admin.DefaultAdminController.GetAdminList))
		adminRouter.GET("/login/info", apiutil.Format(admin.DefaultAdminController.GetLoginAdminInfo))

		pointRouter := adminRouter.Group("/point")
		{
			pointRouter.GET("/record/list", apiutil.Format(admin.DefaultPointController.GetPointRecordPageList))
			pointRouter.GET("/record/list/export", apiutil.Format(admin.DefaultPointController.ExportPointRecordList))
			pointRouter.GET("/type/list", apiutil.Format(admin.DefaultPointController.GetPointTypeList))

			pointRouter.POST("/adjust", apiutil.Format(admin.DefaultPointController.AdjustUserPoint))
			pointRouter.GET("/adjust/list", apiutil.Format(admin.DefaultPointController.GetAdjustRecordPageList))
			pointRouter.GET("/adjust/type/list", apiutil.Format(admin.DefaultPointController.GetAdjustPointTransactionTypeList))
		}

		fileExportRouter := adminRouter.Group("/file-export")
		{
			fileExportRouter.GET("/list", apiutil.Format(admin.DefaultFileExportController.GetFileExportPageList))
			fileExportRouter.GET("/options", apiutil.FormatInterface(admin.DefaultFileExportController.GetFileExportOptions))
		}

		//后台接口-文章管理
		articleRouter := adminRouter.Group("/topic")
		{
			articleRouter.GET("/list", apiutil.Format(admin.DefaultTopicController.List))            //文章列表
			articleRouter.GET("/detail", apiutil.Format(admin.DefaultTopicController.Detail))        //文章详情
			articleRouter.POST("/update", apiutil.Format(admin.DefaultTopicController.Update))       //文章详情
			articleRouter.POST("/delete", apiutil.Format(admin.DefaultTopicController.Delete))       //文章删除
			articleRouter.POST("/review", apiutil.Format(admin.DefaultTopicController.Review))       //审核
			articleRouter.POST("/recommend", apiutil.Format(admin.DefaultTopicController.Recommend)) //推荐
		}
		//后台接口-tag管理
		tagRouter := adminRouter.Group("/tag")
		{
			tagRouter.GET("/list", apiutil.Format(admin.DefaultTagController.List))      //话题列表
			tagRouter.GET("/detail", apiutil.Format(admin.DefaultTagController.Detail))  //话题详情
			tagRouter.POST("/update", apiutil.Format(admin.DefaultTagController.Update)) //更新话题
			tagRouter.POST("/delete", apiutil.Format(admin.DefaultTagController.Delete)) //删除话题
		}
		//后台接口-用户管理
		userRouter := adminRouter.Group("/user")
		{
			userRouter.GET("/list")
			userRouter.GET("/detail")
			userRouter.POST("/change-status")
		}
	}
	adminRouter.POST("/user/list", apiutil.Format(admin.GetUserPageListBy))
	//adminRouter.POST("/user/risk", apiutil.Format(admin.UpdateUserRisk))
}
