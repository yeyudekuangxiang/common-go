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
		adminRouter.GET("/constant", apiutil.Format(admin.DefaultConstantController.List))
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
		channelRouter := adminRouter.Group("/channel")
		{
			channelRouter.POST("/create", apiutil.Format(admin.DefaultUserChannelController.Create))
			channelRouter.POST("/update", apiutil.Format(admin.DefaultUserChannelController.UpdateByCid))
			channelRouter.GET("/list", apiutil.Format(admin.DefaultUserChannelController.GetPageList))
		}
		bannerRouter := adminRouter.Group("/banner")
		{
			bannerRouter.POST("/create", apiutil.Format(admin.DefaultBannerController.Create))
			bannerRouter.POST("/update", apiutil.Format(admin.DefaultBannerController.Update))
			bannerRouter.GET("/list", apiutil.Format(admin.DefaultBannerController.GetPageList))

		}
	}
}
