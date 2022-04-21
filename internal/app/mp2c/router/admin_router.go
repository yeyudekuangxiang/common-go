package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/admin"
	"mio/internal/app/mp2c/middleware"
	"mio/internal/pkg/util"
)

func adminRouter(router *gin.Engine) {

	router.POST("/admin/login", util.Format(admin.DefaultAdminController.Login))
	adminRouter := router.Group("/admin")
	adminRouter.Use(middleware.AuthAdmin())
	{

		adminRouter.GET("/info/list", util.Format(admin.DefaultAdminController.GetAdminList))
		adminRouter.GET("/login/info", util.Format(admin.DefaultAdminController.GetLoginAdminInfo))

		pointRouter := adminRouter.Group("/point")
		{
			pointRouter.GET("/record/list", util.Format(admin.DefaultPointController.GetPointRecordPageList))
			pointRouter.GET("/record/list/export", util.Format(admin.DefaultPointController.ExportPointRecordList))
			pointRouter.GET("/type/list", util.Format(admin.DefaultPointController.GetPointTypeList))

			pointRouter.POST("/adjust", util.Format(admin.DefaultPointController.AdjustUserPoint))
			pointRouter.GET("/adjust/list", util.Format(admin.DefaultPointController.GetAdjustRecordPageList))
			pointRouter.GET("/adjust/type/list", util.Format(admin.DefaultPointController.GetAdjustPointTransactionTypeList))
		}

		fileExportRouter := adminRouter.Group("/file-export")
		{
			fileExportRouter.GET("/list", util.Format(admin.DefaultFileExportController.GetFileExportPageList))
			fileExportRouter.GET("/options", util.FormatInterface(admin.DefaultFileExportController.GetFileExportOptions))
		}
	}
}
