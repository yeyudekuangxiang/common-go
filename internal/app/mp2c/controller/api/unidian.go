package api

import (
	"github.com/gin-gonic/gin"
)

var DefaultUnidianController = UnidianController{}

type UnidianController struct {
}

func (UnidianController) Callback(c *gin.Context) {
	c.String(200, "success")
	return
}
