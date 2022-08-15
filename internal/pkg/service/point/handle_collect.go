package point

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func (c *defaultClientHandle) coffeeCup() error {
	c.withMessage(fmt.Sprintf("coffeeCup=%v", c.clientHandle.additionInfo))
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
	c.withMessage(fmt.Sprintf("bikeRide=%v", c.clientHandle.additionInfo))
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
	c.withMessage(fmt.Sprintf("powerReplace=%v", c.clientHandle.additionInfo))
	_, err := c.incPoint(c.clientHandle.point)
	if err != nil {
		return err
	}
	err = c.saveRecord()
	if err != nil {
		return err
	}
	//计算减碳
	fromString, err := decimal.NewFromString("22.2")
	if err != nil {
		return err
	}
	c.clientHandle.identifyImg["co2"] = fromString.Mul(decimal.NewFromFloat(511)).String()
	delete(c.clientHandle.identifyImg, "orderId")
	return nil
}

func (c *defaultClientHandle) invite() error {
	c.withMessage(fmt.Sprintf("invite=%v", c.clientHandle.additionInfo))
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
	c.withMessage(fmt.Sprintf("article=%v", c.clientHandle.additionInfo))
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
	c.withMessage(fmt.Sprintf("fastElectricity=%v", c.clientHandle.additionInfo))
	//获取充电数 1度电 = 10 积分

	return nil
}
