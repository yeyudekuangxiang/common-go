package point

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

func (c *defaultClientHandle) coffeeCup() error {
	c.withMessage(fmt.Sprintf("coffeeCup=%v", c.clientHandle.message))
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

func (c *defaultClientHandle) bikeRide() error {
	c.withMessage(fmt.Sprintf("bikeRide=%v", c.clientHandle.message))
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

func (c *defaultClientHandle) powerReplace() error {
	c.withMessage(fmt.Sprintf("powerReplace=%v", c.clientHandle.message))
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
	//c.clientHandle.identifyImg 里只有kwh和orderId;需要返回本次积分 本次充电度数 本次减碳 今日累计获得
	data, _, err := c.getTodayData()
	if err != nil {
		return err
	}
	var todayPoint int64
	for _, item := range data {
		todayPoint += item["value"].(int64)
	}
	//计算减碳
	m["co2"] = fromString.Mul(decimal.NewFromFloat(511)).String()
	m["point"] = strconv.FormatInt(c.clientHandle.point, 10)
	m["todayPoint"] = strconv.FormatInt(todayPoint, 10)
	delete(c.clientHandle.identifyImg, "orderId")
	return nil
}

func (c *defaultClientHandle) invite() error {
	c.withMessage(fmt.Sprintf("invite=%v", c.clientHandle.message))
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

func (c *defaultClientHandle) fastElectricity() error {
	c.withMessage(fmt.Sprintf("fastElectricity=%v", c.clientHandle.message))
	//获取充电数 1度电 = 10 积分

	return nil
}
