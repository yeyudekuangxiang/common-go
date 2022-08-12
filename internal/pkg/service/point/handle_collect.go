package point

import "fmt"

func (c *defaultClientHandle) coffeeCup() error {
	c.WithMessage(fmt.Sprintf("coffeeCup=%v", c.clientHandle.additionInfo))
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
	c.WithMessage(fmt.Sprintf("bikeRide=%v", c.clientHandle.additionInfo))
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
	c.WithMessage(fmt.Sprintf("powerReplace=%v", c.clientHandle.additionInfo))
	c.additional.orderId = "t"
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

func (c *defaultClientHandle) invite() error {
	c.WithMessage(fmt.Sprintf("invite=%v", c.clientHandle.additionInfo))
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
	c.WithMessage(fmt.Sprintf("article=%v", c.clientHandle.additionInfo))
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
	c.WithMessage(fmt.Sprintf("fastElectricity=%v", c.clientHandle.additionInfo))
	//获取充电数 1度电 = 10 积分

	return nil
}
