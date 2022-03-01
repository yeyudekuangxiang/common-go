package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/util"
	"mio/service"
)

var DefaultTopicController = TopicController{}

type TopicController struct {
}

func (TopicController) List(c *gin.Context) (gin.H, error) {
	topicTagId := util.PostForm(c, "topicTagId")
	//if err != nil {
	//	return gin.H{
	//		"data": nil,
	//	}, err
	//}
	//
	//page, err := util.PostFormInt(c, "page")
	//if err != nil {
	//	return gin.H{
	//		"data": nil,
	//	}, err
	//}

	list := service.DefaultTopicService.List(util.NewSqlCnd2().EqByReq("topic_tag_id", topicTagId).Desc("sort"))

	return gin.H{
		"data": list,
	}, nil

}
