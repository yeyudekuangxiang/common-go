package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/util"
)

var DefaultChargeController = ChargeController{}

type ChargeController struct {
}

func (ChargeController) Push(c *gin.Context) (gin.H, error) {
	form := GetChargeForm{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
