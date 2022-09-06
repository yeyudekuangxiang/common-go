package point

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"strconv"
	"strings"
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

//oola旧物回收before弹框数据
func (c *defaultClientHandle) oolaRecyclePageData() (map[string]interface{}, error) {
	types := []entity.PointTransactionType{
		entity.POINT_RECYCLING_CLOTHING,
		entity.POINT_RECYCLING_DIGITAL,
		entity.POINT_RECYCLING_APPLIANCE,
		entity.POINT_RECYCLING_BOOK,
	}
	openIds := []string{c.clientHandle.OpenId}
	result, count, err := c.getTodayData(openIds, types)
	if err != nil {
		return nil, err
	}
	//减碳量计算
	var co2, co2Count int64
	res := make(map[string]interface{}, 0)
	//co2
	for _, item := range result {
		additionalInfo := strings.Split(string(item["additional_info"].(entity.AdditionalInfo)), "#")
		co2, _ = strconv.ParseInt(additionalInfo[1], 10, 64)
		co2Count += co2
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2Count
	return res, nil
}

//飞蚂蚁
func (c *defaultClientHandle) fmyRecyclePageData() (map[string]interface{}, error) {
	types := []entity.PointTransactionType{
		entity.POINT_FMY_RECYCLING_CLOTHING,
	}
	openIds := []string{c.clientHandle.OpenId}
	result, count, err := c.getTodayData(openIds, types)
	if err != nil {
		return nil, err
	}
	//减碳量计算
	var co2, co2Count int64
	res := make(map[string]interface{}, 0)
	//co2
	for _, item := range result {
		additionalInfo := strings.Split(string(item["additional_info"].(entity.AdditionalInfo)), "#")
		co2, _ = strconv.ParseInt(additionalInfo[1], 10, 64)
		co2Count += co2
	}
	//返回数据
	res["count"] = count
	res["co2"] = co2Count
	return res, nil
}

//环保减塑
func (c *defaultClientHandle) reducePlasticPageData() (map[string]interface{}, error) {
	types := []entity.PointTransactionType{c.clientHandle.Type}
	openIds := []string{c.clientHandle.OpenId}
	_, count, err := c.getTodayData(openIds, types)
	if err != nil {
		return nil, err
	}
	//返回数据
	res := make(map[string]interface{}, 0)
	res["count"] = count
	res["co2"] = 55 * count
	return res, nil
}

//快电
func (c *defaultClientHandle) fastElectricityPageData() (map[string]interface{}, error) {
	openIds := []string{c.clientHandle.OpenId}
	types := []entity.PointTransactionType{entity.POINT_FAST_ELECTRICITY}
	result, count, err := c.getTodayData(openIds, types)
	if err != nil {
		return nil, err
	}
	//减碳量计算
	res := make(map[string]interface{}, 0)
	//kwh
	var kwhTotal int64
	for _, item := range result {
		kwhTotal += item["value"].(int64) / 10
	}
	//返回数据
	//res["total"] = kwhTotal.String()
	res["count"] = count
	res["co2"] = kwhTotal * 511
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
