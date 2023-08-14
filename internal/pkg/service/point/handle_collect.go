package point

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"regexp"
	"strconv"
)

//奥动 电车换电
func (c *defaultClientHandle) powerReplace() error {

	//c.clientHandle.message
	//c.clientHandle.identifyImg 里只有kwh和orderId;
	m := c.clientHandle.identifyImg
	//检查orderId
	err := c.checkOrderId(m["orderId"])
	if err != nil {
		return err
	}
	re := regexp.MustCompile("\\D") // 匹配非数字字符
	result := re.ReplaceAllString(m["kwh"], "")
	//转换成int
	fromString, err := decimal.NewFromString(result)
	if err != nil {
		return err
	}

	//本次积分
	point := fromString.Mul(decimal.NewFromInt(c.clientHandle.point)).IntPart()
	//检查积分
	if err = c.checkMaxPoint(c.clientHandle.maxPoint, point); err != nil {
		return err
	}
	//更新积分
	err = c.ctx.DB.Transaction(func(tx *gorm.DB) error {
		_, err = c.incPoint(c.clientHandle.point)
		if err != nil {
			return err
		}
		//保存记录
		err = c.saveRecord()
		if err != nil {
			return err
		}
		//更新次数
		err = c.changeLimit()
		if err != nil {
			return err
		}
		kwh, _ := fromString.Float64()
		//发碳量
		_, _ = service.NewCarbonTransactionService(context.NewMioContext()).Create(api_types.CreateCarbonTransactionDto{
			OpenId:  c.clientHandle.OpenId,
			UserId:  0,
			Type:    entity.CarbonTransactionType(c.clientHandle.Type),
			Value:   kwh,
			Info:    m["orderId"] + "#" + m["kwh"] + "#" + c.clientHandle.bizId + "#" + c.clientHandle.ImgUrl,
			AdminId: 0,
			Ip:      "",
			BizId:   c.clientHandle.bizId,
		})
		return nil
	})

	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return err
	}
	data, _, err := c.getPointTodayData([]string{c.clientHandle.OpenId}, []entity.PointTransactionType{c.clientHandle.Type})
	if err != nil {
		return err
	}
	var todayPoint int64
	for _, item := range data {
		todayPoint += item["value"].(int64)
	}
	//计算减碳 需要返回本次积分 本次充电度数 本次减碳 今日累计获得
	m["co2"] = fromString.Mul(decimal.NewFromFloat(511)).String()

	m["point"] = strconv.FormatInt(c.clientHandle.point, 10)
	m["todayPoint"] = strconv.FormatInt(todayPoint, 10)
	delete(c.clientHandle.identifyImg, "orderId")
	return nil
}

//发贴获取积分
func (c *defaultClientHandle) article() error {
	_, err := c.incPoint(c.clientHandle.point)
	if err != nil {
		return err
	}
	err = c.saveRecord()
	if err != nil {
		return err
	}
	return nil
}
