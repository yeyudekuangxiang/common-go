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
	}
	adminRouter.POST("/user/list", apiutil.Format(admin.GetUserPageListBy))
	//adminRouter.POST("/user/risk", apiutil.Format(admin.UpdateUserRisk))
}
