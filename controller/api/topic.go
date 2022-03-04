package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple"
	"mio/internal/util"
	"mio/service"
)

var DefaultTopicController = TopicController{}

type TopicController struct {
}

func (TopicController) List(c *gin.Context) (gin.H, error) {
	form := GetTopicListForm{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}

	sqlCnd := &simple.SqlCnd{}
	if form.TopicTagId != nil {
		sqlCnd.Where("topic_tag_id = ?", form.TopicTagId)
	}
	if form.ID != nil {
		sqlCnd.Where("id = ?", form.ID)
	}
	sqlCnd.Desc("sort")
	list := service.DefaultTopicService.List(sqlCnd)
	return gin.H{
		"list": list,
	}, nil
}
func (TopicController) GetShareWeappQrCode(c *gin.Context) (gin.H, error) {
	form := GetWeappQrCodeFrom{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}
	buffers, contType, err := service.DefaultTopicService.GetShareWeappQrCode(form.OpenId, form.TopicId)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"buffers":     buffers,
		"contentType": contType,
	}, nil
}
