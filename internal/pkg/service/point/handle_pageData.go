package point

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

//电池换电 奥动
func (c *defaultClientHandle) powerReplacePageData() (map[string]interface{}, error) {
	result, count, err := c.getTodayData()
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
		fmt.Printf("%v", s)
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
	_, _, err := c.getTodayData()
	if err != nil {
		return nil, err
	}
	//减碳量计算
	//if _, ok := data["total"]; ok {
	//	data["co2"] = data["total"] * 511
	//}
	return nil, nil
}

func (c *defaultClientHandle) invitePageData() (map[string]interface{}, error) {
	_, _, err := c.getTodayData()
	if err != nil {
		return nil, err
	}
	//减碳量计算
	//if _, ok := data["total"]; ok {
	//	data["co2"] = data["total"] * 511
	//}
	return nil, nil
}

func (c *defaultClientHandle) bikeRidePageData() (map[string]interface{}, error) {
	_, _, err := c.getTodayData()
	if err != nil {
		return nil, err
	}
	//减碳量计算
	//if _, ok := data["total"]; ok {
	//	data["co2"] = data["total"] * 511
	//}
	return nil, nil
}

func (c *defaultClientHandle) articlePageData() (map[string]interface{}, error) {
	_, _, err := c.getTodayData()
	if err != nil {
		return nil, err
	}
	//减碳量计算
	//if _, ok := data["total"]; ok {
	//	data["co2"] = data["total"] * 511
	//}
	return nil, nil
}

func (c *defaultClientHandle) fastElectricityPageData() (map[string]interface{}, error) {
	_, _, err := c.getTodayData()
	if err != nil {
		return nil, err
	}
	//减碳量计算
	//if _, ok := data["total"]; ok {
	//	data["co2"] = data["total"] * 511
	//}
	return nil, nil
}

func (c *defaultClientHandle) getTodayData() ([]map[string]interface{}, int64, error) {
	return repository.NewPointTransactionRepository(c.ctx).CountByToday(repository.GetPointTransactionCountBy{
		OpenIds: []string{c.clientHandle.OpenId},
		Type:    entity.PointTransactionType(c.clientHandle.Type),
	})
}
