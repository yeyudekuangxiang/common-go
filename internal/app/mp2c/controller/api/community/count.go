package community

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/util/apiutil"
)

/*
	酷喵圈计数统计接口
*/

var DefaultCountController = CountController{}

type CountController struct {
}

func (ctr CountController) topicViews(c *gin.Context) (gin.H, error) {
	form := topicCountRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	return nil, nil
}
