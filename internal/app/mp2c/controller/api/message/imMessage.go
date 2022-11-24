package message

import "github.com/gin-gonic/gin"

var DefaultIMMessageController = IMMsgController{}

type IMMsgController struct {
}

func (ctr IMMsgController) Send(c *gin.Context) (gin.H, error) {

	return nil, nil
}

func (ctr IMMsgController) GetByFriend(c *gin.Context) (gin.H, error) {
	return nil, nil
}

func (ctr IMMsgController) BindFriend(c *gin.Context) (gin.H, error) {

	return nil, nil
}
