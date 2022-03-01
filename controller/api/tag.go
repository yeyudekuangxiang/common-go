package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/util"
	"mio/service"
)

var DefaultTagController = TagController{}

type TagController struct {
}

func (TagController) List(c *gin.Context) (gin.H, error) {
	list := service.DefaultTagService.List(util.NewSqlCnd2().Desc("sort"))

	return gin.H{
		"data": list,
	}, nil

}
