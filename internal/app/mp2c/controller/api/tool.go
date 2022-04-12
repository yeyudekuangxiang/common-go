package api

import (
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"mio/internal/pkg/util"
)

var DefaultToolController = ToolController{}

type ToolController struct {
}

func (ToolController) GetQrcode(c *gin.Context) (gin.H, error) {
	form := CreateQrcodeForm{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}
	var png []byte
	png, err := qrcode.Encode(form.Src, qrcode.Medium, 256)
	return gin.H{
		"img":  png,
		"type": "image/png",
	}, err

}
