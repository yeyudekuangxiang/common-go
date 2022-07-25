package admin

import "github.com/gin-gonic/gin"

var DefaultTagController = TagController{}

type TagController struct {
}

func (ctr *TagController) List(context *gin.Context) (gin.H, error) {
	return nil, nil
}

func (ctr *TagController) Detail(context *gin.Context) (gin.H, error) {
	return nil, nil
}

func (ctr *TagController) Update(context *gin.Context) (gin.H, error) {
	return nil, nil

}

func (ctr *TagController) Delete(context *gin.Context) (gin.H, error) {
	return nil, nil

}
