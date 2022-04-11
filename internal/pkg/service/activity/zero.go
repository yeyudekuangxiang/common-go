package activity

import (
	"errors"
	"fmt"
	"mio/internal/pkg/service"
	"time"
)

var (
	ZeroActivityStartTime, _ = time.Parse("2006-01-02 15:04:05", "2022-04-11 14:08:00")
)
var DefaultZeroService = ZeroService{}

type ZeroService struct {
}

func (srv ZeroService) AutoLogin(userId int64) (string, error) {
	userInfo, err := service.DefaultUserService.GetUserById(userId)
	if err != nil {
		return "", err
	}
	if userInfo.ID == 0 {
		return "", errors.New("未查询到用户信息")
	}
	isNewUser := 0
	if userInfo.Time.After(ZeroActivityStartTime) {
		isNewUser = 1
	}
	fmt.Println("asa", fmt.Sprintf("nickname=%snewUser=%d", userInfo.Nickname, isNewUser))
	return service.DefaultDuiBaService.AutoLoginOpenId(service.AutoLoginOpenIdParam{
		UserId:  userId,
		OpenId:  userInfo.OpenId,
		Path:    "https://88543.activity-12.m.duiba.com.cn/wechat/access?apk=ngiGp48EcRUC9TjpXEYxdSSJhim&dbredirect=https%3A%2F%2F88543.activity-12.m.duiba.com.cn%2Faaw%2Fhaggle%2Findex%3FopId%3D194935804526281%26dbnewopen%26newChannelType%3D3",
		DCustom: "newUser12=12",
	})
}
