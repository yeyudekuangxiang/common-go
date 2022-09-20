package message

import (
	"github.com/gin-gonic/gin"
	"mio/config"
	messageSrv "mio/internal/pkg/service/message"
)

var DefaultMessageController = MessageController{}

type MessageController struct {
}

func (MessageController) SendMessage(c *gin.Context) (gin.H, error) {
	b := messageSrv.MiniChangePointTemplate{
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

func (MessageController) GetTemplateId(c *gin.Context) (gin.H, error) {
	var TemplateIds []string
	TemplateIds = append(TemplateIds, config.MessageTemplateIds.ChangePoint)
	TemplateIds = append(TemplateIds, config.MessageTemplateIds.ChangePoint)
	return gin.H{
		"templateIds": TemplateIds,
	}, nil
}
