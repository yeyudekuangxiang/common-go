package open

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/util/apiutil"
)

var DefaultJhxController = JhxController{}

type JhxController struct {
}

func (ctr JhxController) BusTicketNotify(ctx *gin.Context) (gin.H, error) {
	form := jhxUseCodeFrom{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	return nil, nil
}
