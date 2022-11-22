package message

import (
	"github.com/gin-gonic/gin"
	"mio/config"
	"mio/internal/pkg/core/context"
	messageSrv "mio/internal/pkg/service/message"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"strings"
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

//获取模板
func (ctr MsgController) GetTemplateId(c *gin.Context) (gin.H, error) {
	form := MessageGetTemplateIdForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	service := messageSrv.MessageService{}
	return gin.H{
		"templateIds": service.GetTemplateId(user.OpenId, form.Scene),
	}, nil
}

//获取消息
func (ctr MsgController) GetWebMessage(c *gin.Context) (gin.H, error) {
	form := WebMessageRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	types := strings.Split(strings.Trim(form.Types, ","), ",")

	if types[0] == "" {
		return nil, errno.ErrCommon.WithMessage("参数错误")
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	user := apiutil.GetAuthUser(c)

	messageService := messageSrv.NewWebMessageService(ctx)
	msgList, total, err := messageService.GetMessage(messageSrv.GetWebMessage{
		UserId: user.ID,
		Status: form.Status,
		Types:  types,
		Limit:  form.Limit(),
		Offset: form.Offset(),
	})

	if err != nil {
		return nil, err
	}

	return gin.H{
		"list":     msgList,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

//获取消息数量
func (ctr MsgController) GetWebMessageCount(c *gin.Context) (gin.H, error) {
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	user := apiutil.GetAuthUser(c)

	messageService := messageSrv.NewWebMessageService(ctx)
	resp, err := messageService.GetMessageCount(messageSrv.GetWebMessageCount{
		RecId: user.ID,
	})
	if err != nil {
		return nil, err
	}

	return gin.H{
		"total":            resp.Total,
		"exchangeMsgTotal": resp.ExchangeMsgTotal,
		"systemMsgTotal":   resp.SystemMsgTotal,
	}, nil
}

//设置已读
func (ctr MsgController) SetHaveReadWebMessage(c *gin.Context) (gin.H, error) {
	form := HaveReadWebMessageRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	var msgIds []string
	if form.MsgIds != "" {
		ids := strings.Split(strings.Trim(form.MsgIds, ","), ",")
		if ids[0] != "" {
			msgIds = ids
		}
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	user := apiutil.GetAuthUser(c)

	messageService := messageSrv.NewWebMessageService(ctx)
	err := messageService.SetHaveRead(messageSrv.SetHaveReadMessage{
		RecId:  user.ID,
		MsgIds: msgIds,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//消息跳转
func (ctr MsgController) TurnWebMessage(c *gin.Context) (gin.H, error) {
	//form := TurnWebMessageRequest{}
	//if err := apiutil.BindForm(c, &form); err != nil {
	//	return nil, err
	//}
	//
	//ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	//user := apiutil.GetAuthUser(c)
	//
	return nil, nil
}
