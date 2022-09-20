package message

import (
	"github.com/gin-gonic/gin"
	messageSrv "mio/internal/pkg/service/message"
)

var DefaultMessageController = MessageController{}

type MessageController struct {
}

func (MessageController) SendMessage(c *gin.Context) (gin.H, error) {
	b := messageSrv.MiniPointSendTemplate{
		Point:    "222",
		Source:   "222",
		Time:     "2022年9月10日",
		AllPoint: "222",
	}
	service := messageSrv.MessageService{}
	code, err := service.SendMiniSubMessage("oy_BA5IGl1JgkJKbD14wq_-Yorqw", "index", b)
	return gin.H{
		"code": code,
		"err":  err,
	}, nil
}
