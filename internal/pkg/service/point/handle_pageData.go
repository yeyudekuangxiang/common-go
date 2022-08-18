package point

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
)

//电池换电 奥动 before 弹框数据
func (c *defaultClientHandle) powerReplacePageData() (map[string]interface{}, error) {
	types := []entity.PointTransactionType{c.clientHandle.Type}
	openIds := []string{c.clientHandle.OpenId}
	result, count, err := c.getTodayData(openIds, types)
	if err != nil {
		return nil, err
	}
	//减碳量计算
	m := make(map[string]interface{}, 0)
	res := make(map[string]interface{}, 0)
	//kwh
	kwhTotal := decimal.NewFromFloat(0)
	for _, item := range result {
		s := item["additional_info"].(entity.AdditionalInfo)
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			return nil, err
		}
		fromString, _ := decimal.NewFromString(m["kwh"].(string))
		kwhTotal = fromString.Add(kwhTotal)
	}
	//返回数据
	//res["total"] = kwhTotal.String()
	res["count"] = count
	res["co2"] = kwhTotal.Mul(decimal.NewFromFloat(511)).String()
	return res, nil
}

func (c *defaultClientHandle) coffeeCupPageData() (map[string]interface{}, error) {
	//_, _, err := c.getTodayData()
	//if err != nil {
	//	return nil, err
	//}
	//减碳量计算
	//if _, ok := data["total"]; ok {
	//	data["co2"] = data["total"] * 511
	//}
	return nil, nil
}

func (c *defaultClientHandle) bikeRidePageData() (map[string]interface{}, error) {
	//_, _, err := c.getTodayData()
	//if err != nil {
	//	return nil, err
	//}
	//减碳量计算
	//if _, ok := data["total"]; ok {
	//	data["co2"] = data["total"] * 511
	//}
	return nil, nil
}

//oola旧物回收before弹框数据
func (c *defaultClientHandle) oolaRecyclePageData() (map[string]interface{}, error) {
	types := []entity.PointTransactionType{
		entity.POINT_OOLA_RECYCLING_CLOTHING,
		entity.POINT_OOLA_RECYCLING_DIGITAL,
		entity.POINT_OOLA_RECYCLING_APPLIANCE,
		entity.POINT_OOLA_RECYCLING_BOOK,
	}
	openIds := []string{c.clientHandle.OpenId}
	result, count, err := c.getTodayData(openIds, types)
	if err != nil {
		return nil, err
	}
	//减碳量计算
	var co2 int64
	res := make(map[string]interface{}, 0)
	//kwh
	for _, item := range result {
		temp := service.DefaultRecycleService.GetCo2(entity.PointTransactionType(item["type"].(string)), item["value"].(int64))
		co2 += temp
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2
	return res, nil
}

func (c *defaultClientHandle) getTodayData(openIds []string, types []entity.PointTransactionType) ([]map[string]interface{}, int64, error) {
	return repository.NewPointTransactionRepository(c.ctx).CountByToday(repository.GetPointTransactionCountBy{
		OpenIds: openIds,
		Types:   types,
	})
}

func (c *defaultClientHandle) getMonthData(openIds []string, types []entity.PointTransactionType) ([]map[string]interface{}, int64, error) {
	return repository.NewPointTransactionRepository(c.ctx).CountByMonth(repository.GetPointTransactionCountBy{
		OpenIds: openIds,
		Types:   types,
	})
}
