package message

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	messageSrv "mio/internal/pkg/service/message"
	"mio/internal/pkg/util/apiutil"
)

var DefaultMessageController = MessageController{}

type MessageController struct {
}

func (MessageController) SendMessage(c *gin.Context) (gin.H, error) {
	b := messageSrv.MiniChangePointTemplate{
		Point:    2222,
		Source:   "222",
		Time:     "2022年9月10日",
		AllPoint: 333,
	}
	service := messageSrv.MessageService{}
	code, err := service.SendMiniSubMessage("oy_BA5IGl1JgkJKbD14wq_-Yorqw", "index", b)
	return gin.H{
		"code": code,
		"err":  err,
	}, nil
}

func (MessageController) GetTemplateId(c *gin.Context) (gin.H, error) {
	form := api_types.MessageGetTemplateIdForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	service := messageSrv.MessageService{}
	return gin.H{
		"templateIds": service.GetTemplateId(user.OpenId, form.Scene),
	}, nil
}
