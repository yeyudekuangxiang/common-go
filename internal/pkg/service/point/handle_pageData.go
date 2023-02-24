package point

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
)

//电池换电 奥动 before 弹框数据
func (c *defaultClientHandle) powerReplacePageData() (map[string]interface{}, error) {
	types := c.getCarbonType(string(c.clientHandle.Type))
	result, count, err := c.getCarbonDayData(c.clientHandle.OpenId, types)
	if err != nil {
		return nil, err
	}
	var co2 float64
	res := make(map[string]interface{}, 0)
	for _, item := range result {
		co2 += item.Value
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2
	return res, nil
}

//oola旧物回收before弹框数据
func (c *defaultClientHandle) oolaRecyclePageData() (map[string]interface{}, error) {
	types := c.getCarbonType(string(c.clientHandle.Type))
	result, count, err := c.getCarbonDayData(c.clientHandle.OpenId, types)
	if err != nil {
		return nil, err
	}
	//减碳量计算
	var co2 float64
	res := make(map[string]interface{}, 0)
	//co2
	for _, item := range result {
		co2 += item.Value
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2
	return res, nil
}

//飞蚂蚁
func (c *defaultClientHandle) fmyRecyclePageData() (map[string]interface{}, error) {
	types := c.getCarbonType(string(c.clientHandle.Type))
	result, count, err := c.getCarbonDayData(c.clientHandle.OpenId, types)
	if err != nil {
		return nil, err
	}
	//减碳量计算
	var co2 float64
	res := make(map[string]interface{}, 0)
	//co2
	for _, item := range result {
		co2 += item.Value
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2
	return res, nil
}

//环保减塑
func (c *defaultClientHandle) reducePlasticPageData() (map[string]interface{}, error) {
	types := c.getCarbonType(string(c.clientHandle.Type))
	result, count, err := c.getCarbonDayData(c.clientHandle.OpenId, types)
	if err != nil {
		return nil, err
	}
	var co2 float64
	res := make(map[string]interface{}, 0)

	for _, item := range result {
		co2 += item.Value
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2
	return res, nil
}

//toCharge 快电、云快充
func (c *defaultClientHandle) toChargePageData() (map[string]interface{}, error) {
	types := c.getCarbonType(string(c.clientHandle.Type))
	result, count, err := c.getCarbonDayData(c.clientHandle.OpenId, types)
	if err != nil {
		return nil, err
	}
	//减碳量计算
	res := make(map[string]interface{}, 0)
	var co2 float64

	for _, item := range result {
		co2 += item.Value
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2
	return res, nil
}

//金华行
func (c *defaultClientHandle) jhxPageData() (map[string]interface{}, error) {
	types := c.getCarbonType(string(c.clientHandle.Type))
	result, count, err := c.getCarbonDayData(c.clientHandle.OpenId, types)
	if err != nil {
		return nil, err
	}
	var co2 float64
	res := make(map[string]interface{}, 0)
	for _, item := range result {
		co2 += item.Value
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2
	return res, nil
}

//亿通行
func (c *defaultClientHandle) ytxPageData() (map[string]interface{}, error) {
	types := c.getCarbonType(string(c.clientHandle.Type))
	result, count, err := c.getCarbonDayData(c.clientHandle.OpenId, types)
	if err != nil {
		return nil, err
	}
	var co2 float64
	res := make(map[string]interface{}, 0)
	for _, item := range result {
		co2 += item.Value
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2
	return res, nil
}

func (c *defaultClientHandle) ahsRecyclePageData() (map[string]interface{}, error) {
	types := c.getCarbonType(string(c.clientHandle.Type))
	result, count, err := c.getCarbonDayData(c.clientHandle.OpenId, types)
	if err != nil {
		return nil, err
	}
	var co2 float64
	res := make(map[string]interface{}, 0)
	for _, item := range result {
		co2 += item.Value
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2
	return res, nil
}

func (c *defaultClientHandle) sshsRecyclePageData() (map[string]interface{}, error) {
	types := c.getCarbonType(string(c.clientHandle.Type))
	result, count, err := c.getCarbonDayData(c.clientHandle.OpenId, types)
	if err != nil {
		return nil, err
	}
	var co2 float64
	res := make(map[string]interface{}, 0)
	for _, item := range result {
		co2 += item.Value
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2
	return res, nil
}

func (c *defaultClientHandle) ddyxRecyclePageData() (map[string]interface{}, error) {
	types := c.getCarbonType(string(c.clientHandle.Type))
	result, count, err := c.getCarbonDayData(c.clientHandle.OpenId, types)
	if err != nil {
		return nil, err
	}
	var co2 float64
	res := make(map[string]interface{}, 0)
	for _, item := range result {
		co2 += item.Value
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2
	return res, nil
}

func (c *defaultClientHandle) getPointTodayData(openIds []string, types []entity.PointTransactionType) ([]map[string]interface{}, int64, error) {
	return repository.NewPointTransactionRepository(c.ctx).CountByToday(repository.GetPointTransactionCountBy{
		OpenIds: openIds,
		Types:   types,
	})
}

func (c *defaultClientHandle) getCarbonDayData(openId string, tps []entity.CarbonTransactionType) ([]entity.CarbonTransaction, int64, error) {
	result, count, err := service.NewCarbonTransactionService(c.ctx).GetTodayCarbonByType(openId, tps)
	if err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return nil, 0, err
	}
	return result, count, nil
}

func (c *defaultClientHandle) getMonthData(openIds []string, types []entity.PointTransactionType) ([]map[string]interface{}, int64, error) {
	return repository.NewPointTransactionRepository(c.ctx).CountByMonth(repository.GetPointTransactionCountBy{
		OpenIds: openIds,
		Types:   types,
	})
}

func (c *defaultClientHandle) getPointType(tp string) []entity.PointTransactionType {
	switch tp {
	case "POWER_REPLACE":
		return []entity.PointTransactionType{entity.POINT_POWER_REPLACE} //换电
	case "FAST_ELECTRICITY":
		return []entity.PointTransactionType{entity.POINT_FAST_ELECTRICITY} //快电
	case "YKC":
		return []entity.PointTransactionType{entity.POINT_YKC} //云快充
	case "REDUCE_PLASTIC":
		return []entity.PointTransactionType{entity.POINT_REDUCE_PLASTIC}
	case "JHX":
		return []entity.PointTransactionType{entity.POINT_JHX}
	case "YTX":
		return []entity.PointTransactionType{entity.POINT_YTX}
	case "OOLA_RECYCLE":
		return []entity.PointTransactionType{entity.POINT_RECYCLING_CLOTHING}
	case "FMY_RECYCLE":
		return []entity.PointTransactionType{entity.POINT_FMY_RECYCLING_CLOTHING}
	case "AHS_RECYCLE":
		return []entity.PointTransactionType{entity.POINT_RECYCLING_AIHUISHOU} //爱回收旧物回收
	case "SSHS_RECYCLE":
		return []entity.PointTransactionType{entity.POINT_RECYCLING_SHISHANGHUISHOU}
	case "DDYX_RECYCLE":
		return []entity.PointTransactionType{entity.POINT_RECYCLING_DANGDANGYIXIA}
	}
	return []entity.PointTransactionType{}
}
func (c *defaultClientHandle) getCarbonType(tp string) []entity.CarbonTransactionType {
	switch tp {
	case "POWER_REPLACE":
		return []entity.CarbonTransactionType{entity.CARBON_POWER_REPLACE} //换电
	case "FAST_ELECTRICITY":
		return []entity.CarbonTransactionType{entity.CARBON_FAST_ELECTRICITY} //快电
	case "YKC":
		return []entity.CarbonTransactionType{entity.CARBON_YKC} //云快充
	case "REDUCE_PLASTIC":
		return []entity.CarbonTransactionType{entity.CARBON_REDUCE_PLASTIC}
	case "JHX":
		return []entity.CarbonTransactionType{entity.CARBON_JHX}
	case "YTX":
		return []entity.CarbonTransactionType{entity.CARBON_YTX}
	case "OOLA_RECYCLE":
		return []entity.CarbonTransactionType{entity.CARBON_RECYCLING_CLOTHING, entity.CARBON_RECYCLING_DIGITAL, entity.CARBON_RECYCLING_APPLIANCE, entity.CARBON_RECYCLING_BOOK}
	case "FMY_RECYCLE":
		return []entity.CarbonTransactionType{entity.CARBON_FMY_RECYCLING_CLOTHING}
	case "AHS_RECYCLE":
		return []entity.CarbonTransactionType{entity.CARBON_RECYCLING_AIHUISHOU} //爱回收旧物回收
	case "SSHS_RECYCLE":
		return []entity.CarbonTransactionType{entity.CARBON_RECYCLING_SHISHANGHUISHOU}
	case "DDYX_RECYCLE":
		return []entity.CarbonTransactionType{entity.CARBON_RECYCLING_DANGDANGYIXIA}
	}
	return []entity.CarbonTransactionType{}
}
