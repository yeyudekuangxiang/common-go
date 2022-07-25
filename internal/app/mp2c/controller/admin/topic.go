package admin

import "github.com/gin-gonic/gin"

var DefaultTopicController = TopicController{}

type TopicController struct {
}

func (ctr *TopicController) List(ctx *gin.Context) (gin.H, error) {
	return nil, nil
}

func (ctr *TopicController) Detail(context *gin.Context) (gin.H, error) {
	return nil, nil
}

func (ctr *TopicController) Update(context *gin.Context) (gin.H, error) {
	return nil, nil

}

func (ctr *TopicController) Delete(context *gin.Context) (gin.H, error) {
	return nil, nil

}

func (ctr *TopicController) Review(context *gin.Context) (gin.H, error) {
	return nil, nil

}

func (ctr *TopicController) Recommend(context *gin.Context) (gin.H, error) {
	return nil, nil

}
