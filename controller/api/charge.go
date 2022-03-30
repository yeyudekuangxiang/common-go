package api

import (
	"github.com/gin-gonic/gin"
)

var DefaultChargeController = ChargeController{}

type ChargeController struct {
}

func (ChargeController) Push(c *gin.Context) (gin.H, error) {
	return gin.H{}, nil
}
