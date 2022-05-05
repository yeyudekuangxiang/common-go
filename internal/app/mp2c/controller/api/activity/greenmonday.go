package activity

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service/activity"
	"mio/internal/pkg/util/apiutil"
)

var DefaultGMController = GMController{}

type GMController struct {
}

func (ctr GMController) GetGMRecord(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	record, err := activity.DefaultGMService.FindOrCreateGMRecord(user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"record": record,
	}, nil
}

func (ctr GMController) ReportInvitationRecord(ctx *gin.Context) (gin.H, error) {
	form := ReportInvitationRecordForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	if form.UserId == 0 {
		return nil, nil
	}
	user := apiutil.GetAuthUser(ctx)

	err := activity.DefaultGMService.AddInvitationRecord(form.UserId, user.ID)
	if err != nil {
		app.Logger.Error(err)
		return nil, errors.New("建立邀请关系失败,请联系管理员")
	}
	return nil, nil
}

func (ctr GMController) ExchangeGift(ctx *gin.Context) (interface{}, error) {
	form := ExchangeGiftForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	order, err := activity.DefaultGMService.Order(user.ID, form.AddressId)
	if err != nil {
		return nil, err
	}
	return order.ShortOrder(), nil
}

func (ctr GMController) AnswerQuestion(ctx *gin.Context) (gin.H, error) {
	form := GMAnswerQuestion{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	record, err := activity.DefaultGMService.AnswerQuestion(activity.AnswerGMQuestionParam{
		UserId:  user.ID,
		Title:   form.Title,
		Answer:  form.Answer,
		IsRight: form.IsRight,
	})
	return gin.H{
		"record": record,
	}, err
}
