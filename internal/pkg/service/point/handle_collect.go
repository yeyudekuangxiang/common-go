package point

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

//奥动 电车换电
func (c *defaultClientHandle) powerReplace() error {
	c.withMessage(fmt.Sprintf("powerReplace=%v", c.clientHandle.message))
	//c.clientHandle.identifyImg 里只有kwh和orderId;
	m := c.clientHandle.identifyImg
	fromString, err := decimal.NewFromString(m["kwh"])
	if err != nil {
		return err
	}
	c.clientHandle.point = fromString.Mul(decimal.NewFromInt(c.clientHandle.point)).IntPart()
	_, err = c.incPoint(c.clientHandle.point)
	if err != nil {
		return err
	}
	err = c.saveRecord()
	if err != nil {
		return err
	}
	data, _, err := c.getTodayData()
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
	c.withMessage(fmt.Sprintf("article=%v", c.clientHandle.message))
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
