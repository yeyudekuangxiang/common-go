package point

import "fmt"

func (c *clientHandle) coffeeCup() error {
	c.WithMessage(fmt.Sprintf("coffeeCup=%v", c.AdditionInfo))
	_, err := c.incPoint(c.Point)
	if err != nil {
		return err
	}
	err = c.saveRecord()
	if err != nil {
		return err
	}

	return nil
}

func (c *clientHandle) bikeRide() error {
	c.WithMessage(fmt.Sprintf("bikeRide=%v", c.AdditionInfo))
	_, err := c.incPoint(c.Point)
	if err != nil {
		return err
	}
	err = c.saveRecord()
	if err != nil {
		return err
	}

	return nil
}

func (c *clientHandle) powerReplace() error {
	c.WithMessage(fmt.Sprintf("powerReplace=%v", c.AdditionInfo))
	c.additional.orderId = "t"
	_, err := c.incPoint(c.Point)
	if err != nil {
		return err
	}
	err = c.saveRecord()
	if err != nil {
		return err
	}

	return nil
}

func (c *clientHandle) invite() error {
	c.WithMessage(fmt.Sprintf("invite=%v", c.AdditionInfo))
	_, err := c.incPoint(c.Point)
	if err != nil {
		return err
	}
	err = c.saveRecord()
	if err != nil {
		return err
	}

	return nil
}

func (c *clientHandle) article() error {
	c.WithMessage(fmt.Sprintf("article=%v", c.AdditionInfo))
	_, err := c.incPoint(c.Point)
	if err != nil {
		return err
	}
	err = c.saveRecord()
	if err != nil {
		return err
	}
	return nil
}

func (c *clientHandle) fastElectricity() error {
	c.WithMessage(fmt.Sprintf("fastElectricity=%v", c.AdditionInfo))
	//获取充电数 1度电 = 10 积分

	return nil
}
