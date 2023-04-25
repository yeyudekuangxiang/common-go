package middleware

import (
	"github.com/gin-gonic/gin"
)

func BanSpeech() gin.HandlerFunc {
	return func(c *gin.Context) {
		//user := apiutil.GetAuthUser(c)
		//authMap := make(map[int]struct{}, 0)
		//for _, item := range user.Auth {
		//	authMap[item] = struct{}{}
		//}
		//method := c.Request.Method
		//strings.Split(strings.Trim(strings.ToLower(c.FullPath()), "/"), "/")
		//fmt.Println(method)
	}
}
