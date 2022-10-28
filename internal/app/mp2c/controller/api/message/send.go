package message

import (
	"github.com/gin-gonic/gin"
	"mio/config"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	messageSrv "mio/internal/pkg/service/message"
	"mio/internal/pkg/util/apiutil"
)

var DefaultMessageController = MsgController{}

type MsgController struct {
}

//todo 给测试用

func (ctr MsgController) SendMessage(c *gin.Context) (gin.H, error) {

	return gin.H{}, nil
	b := messageSrv.MiniSignRemindTemplate{
		ActivityName: "2323",
		Tip:          "232323",
	}
	service := messageSrv.MessageService{}
	code, err := service.SendMiniSubMessage("oy_BA5DGmQBqMeCj_9Eozj8dXhoA", config.MessageJumpUrls.SignRemind, b)
	return gin.H{
		"code": code,
		"err":  err,
	}, nil

	/*	b := messageSrv.MiniChangePointTemplate{
			Point:    2,
			Source:   "222",
			Time:     "2021-09-10",
			AllPoint: 3,
		}
		service := messageSrv.MessageService{}
		code, err := service.SendMiniSubMessage("oy_BA5DGmQBqMeCj_9Eozj8dXhoA", config.MessageJumpUrls.ChangePoint, b)
		return gin.H{
			"code": code,
			"err":  err,
		}, nil*/

	/*	b := messageSrv.MiniOrderDeliverTemplate{
			OrderNo:      "111",
			TrackNo:      "222",
			TrackCompany: "333",
			GoodName:     "1111",
			Tip:          "121212",
		}
		service := messageSrv.MessageService{}
		code, err := service.SendMiniSubMessage("oy_BA5DGmQBqMeCj_9Eozj8dXhoA", config.MessageJumpUrls.OrderDeliver, b)
		return gin.H{
			"code": code,
			"err":  err,
		}, nil*/
}

//todo 给测试用

func (ctr MsgController) SendSign(c *gin.Context) (gin.H, error) {
	return gin.H{}, nil
	service := messageSrv.MessageService{}
	service.SendMessageToSignUser()
	return gin.H{
		"code": 0,
		"err":  nil,
	}, nil
}

func (ctr MsgController) GetTemplateId(c *gin.Context) (gin.H, error) {
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

func (ctr MsgController) GetWebMessage(c *gin.Context) (gin.H, error) {
	form := api_types.WebMessageRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	user := apiutil.GetAuthUser(c)

	messageService := messageSrv.NewWebMessageService(ctx)
	msgList, total, err := messageService.GetMessage(user.ID, form.Status, form.Limit(), form.Offset())
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list":     msgList,
		"total":    total,
		"page":     form.Limit(),
		"pageSize": form.Offset(),
	}, nil
}
