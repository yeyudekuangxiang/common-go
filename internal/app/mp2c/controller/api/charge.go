package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/util/apiutil"
)

var DefaultChargeController = ChargeController{}

type ChargeController struct {
}

func (ChargeController) Push(c *gin.Context) (gin.H, error) {
	form := GetChargeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	fmt.Println(&form)
	return gin.H{}, nil
}
