package point

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
)

func (c *defaultClientHandle) powerReplacePageData() (map[string]int64, error) {
	data, err := c.getPageData()
	if err != nil {
		return nil, err
	}
	//减碳量计算 1 度电 = 511 CO2
	if _, ok := data["total"]; ok {
		data["co2"] = data["total"] * 511
	}
	return data, nil
}

func (c *defaultClientHandle) coffeeCupPageData() (map[string]int64, error) {
	data, err := c.getPageData()
	if err != nil {
		return nil, err
	}
	//减碳量计算
	data["decCO2"] = 0
	//if _, ok := data["total"]; ok {
	//	data["decCO2"] = data["total"] * 511
	//}
	return data, nil
}

func (c *defaultClientHandle) invitePageData() (map[string]int64, error) {
	data, err := c.getPageData()
	if err != nil {
		return nil, err
	}
	//减碳量计算
	data["decCO2"] = 0
	//if _, ok := data["total"]; ok {
	//	data["decCO2"] = data["total"] * 511
	//}
	return data, nil
}

func (c *defaultClientHandle) bikeRidePageData() (map[string]int64, error) {
	data, err := c.getPageData()
	if err != nil {
		return nil, err
	}
	//减碳量计算
	data["decCO2"] = 0
	//if _, ok := data["total"]; ok {
	//	data["decCO2"] = data["total"] * 511
	//}
	return data, nil
}

func (c *defaultClientHandle) articlePageData() (map[string]int64, error) {
	data, err := c.getPageData()
	if err != nil {
		return nil, err
	}
	//减碳量计算
	data["decCO2"] = 0
	//if _, ok := data["total"]; ok {
	//	data["decCO2"] = data["total"] * 511
	//}
	return data, nil
}

func (c *defaultClientHandle) fastElectricityPageData() (map[string]int64, error) {
	data, err := c.getPageData()
	if err != nil {
		return nil, err
	}
	//减碳量计算
	data["decCO2"] = 0
	//if _, ok := data["total"]; ok {
	//	data["decCO2"] = data["total"] * 511
	//}
	return data, nil
}

func (c *defaultClientHandle) getPageData() (map[string]int64, error) {
	result, err := repository.NewPointTransactionRepository(c.ctx).CountByToday(repository.GetPointTransactionCountBy{
		OpenIds: []string{c.clientHandle.OpenId},
		Type:    entity.PointTransactionType(c.clientHandle.Type),
	})
	if err != nil {
		return nil, err
	}
	data := util.MapInterface2int64(result)
	if _, ok := data["total"]; !ok {
		data["total"] = 0
	}
	if _, ok := data["co2"]; !ok {
		data["co2"] = 0
	}
	if _, ok := data["count"]; !ok {
		data["count"] = 0
	}
	return data, nil
}
