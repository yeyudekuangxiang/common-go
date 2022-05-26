package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
	"strconv"
	"time"
)

var DefaultChargeController = ChargeController{}

type ChargeController struct {
}

func (ChargeController) Push(c *gin.Context) (gin.H, error) {
	form := GetChargeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	fmt.Println("xingxingcharge", form)
	//通过手机号查询用户
	userInfo, _ := service.DefaultUserService.GetUserBy(repository.GetUserBy{Mobile: form.Mobile})
	if userInfo.ID <= 0 {
		fmt.Println("xingxingcharge 未找到用户 ", form.Mobile)
	} else {
		//查询今日积分总量
		timeStr := time.Now().Format("2006-01-02")
		key := timeStr + "xingxingcharge" + form.Mobile
		cmd := app.Redis.Get(c, key)
		lastPoint, _ := strconv.Atoi(cmd.Val())
		thisPoint := int(form.TotalPower * 10)
		totalPoint := lastPoint + thisPoint
		fmt.Println("xingxingcharge 累计积分 ", form.Mobile, totalPoint)
		if lastPoint >= 300 {
			fmt.Println("xingxingcharge 充电量已达到上线 ", form.Mobile)
		} else {
			if totalPoint > 300 {
				fmt.Println("xingxingcharge 充电量限制修正 ", form.Mobile, thisPoint, lastPoint)
				thisPoint = 300 - lastPoint
				totalPoint = 300
			}
		}
		app.Redis.Set(c, key, totalPoint, 24*36000*time.Second)

		//加积分
		_, err := service.DefaultPointTransactionService.Create(service.CreatePointTransactionParam{
			OpenId:       userInfo.OpenId,
			Type:         entity.POINT_ECAR,
			Value:        thisPoint,
			AdditionInfo: form.OutTradeNo + "#" + form.Mobile + "#" + form.Ch + "#" + strconv.Itoa(int(form.TotalPower*10)) + "#" + form.Sign,
		})
		if err != nil {
			fmt.Println("xingxingcharge 加积分失败 ", form.Mobile)
		}
	}

	return gin.H{}, nil
}
