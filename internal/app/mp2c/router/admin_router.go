package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/admin"
	"mio/internal/pkg/util"
)

func adminRouter(router *gin.Engine) {
	adminRouter := router.Group("/admin")
	/*adminRouter.Use(middleware.AuthAdmin())*/
	{
		adminRouter.GET("/user", util.Format(admin.DefaultUserController.GetUserInfo))

		pointRouter := adminRouter.Group("/point")
		{
			pointRouter.GET("/record/list", util.Format(admin.DefaultPointController.GetPointRecordPageList))
			pointRouter.GET("/record/list/export", util.Format(admin.DefaultPointController.ExportPointRecordList))
			pointRouter.GET("/type/list", util.Format(admin.DefaultPointController.GetPointTypeList))
		}

		fileExportRouter := adminRouter.Group("/file-export")
		{
			fileExportRouter.GET("/list", util.Format(admin.DefaultFileExportController.GetFileExportPageList))
			fileExportRouter.GET("/options", util.FormatInterface(admin.DefaultFileExportController.GetFileExportOptions))
		}
	}
}
