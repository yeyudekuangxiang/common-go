package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util/encrypt"
	"mio/pkg/errno"
	"sort"
	"strconv"
	"strings"
)

var DefaultRecycleService = RecycleService{}

type RecycleService struct {
}

//旧物回收 积分规则
var pointCollectByRecycle = map[string]int64{
	"衣帽鞋包": 2709, // 人/月
	"书籍课本": 322,  // 人/月
	"数码产品": 1911, // 人/月
	"家电回收": 1840, // 人/月
}

var pointRecycleByRules = map[string]entity.PointTransactionType{
	"衣帽鞋包": entity.POINT_OOLA_RECYCLING_CLOTHING,  // 人/月
	"书籍课本": entity.POINT_OOLA_RECYCLING_BOOK,      // 人/月
	"数码产品": entity.POINT_OOLA_RECYCLING_DIGITAL,   // 人/月
	"家电回收": entity.POINT_OOLA_RECYCLING_APPLIANCE, // 人/月
}

var pointRecycleByQuantity = map[string]int64{
	"手机":    113,
	"平板电脑":  409,
	"手提电脑":  1031,
	"一体机电脑": 1911,
	"冰箱":    384,
	"洗衣机":   690,
	"空调":    205,
	"电视机":   69,
}

var pointRecycleByLimit = map[string]int64{
	"衣帽鞋包":  1,
	"书籍课本":  1,
	"手机":    1,
	"平板电脑":  1,
	"手提电脑":  1,
	"一体机电脑": 1,
	"冰箱":    1,
	"洗衣机":   1,
	"空调":    1,
	"电视机":   1,
}

var Co2RecycleByPoint = map[entity.PointTransactionType]map[int64]int64{
	entity.POINT_OOLA_RECYCLING_CLOTHING: {
		208:  91,
		417:  181,
		834:  362,
		2709: 1178,
	},
	entity.POINT_OOLA_RECYCLING_BOOK: {
		64:  1,
		129: 3,
		193: 4,
		322: 7,
	},
	entity.POINT_OOLA_RECYCLING_DIGITAL: {
		113:  49,
		409:  178,
		1031: 448,
		1911: 831,
	},
	entity.POINT_OOLA_RECYCLING_APPLIANCE: {
		384: 167,
		690: 300,
		205: 89,
		69:  30,
	},
}

func (srv RecycleService) CheckSign(params map[string]interface{}, secret string) error {
	sign := params["sign"]
	delete(params, "sign")
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += params[v].(string)
	}
	//验证签名
	signMd5 := encrypt.Md5(secret + signStr)
	if signMd5 != sign {
		return errno.ErrValidation
	}
	return nil
}

// GetPoint 计算积分和co2
func (srv RecycleService) GetPoint(typeName string, qua string, unit string) (int64, int64) {
	atoi, _ := strconv.Atoi(qua)
	return srv.getPoint(typeName, atoi)
}

func (srv RecycleService) GetCo2(typeName entity.PointTransactionType, point int64) int64 {
	m := Co2RecycleByPoint[typeName]
	if co2, ok := m[point]; ok {
		return co2
	}
	return 0
}

func (srv RecycleService) GetMaxPointByMonth(typeName string) int64 {
	if point, ok := pointCollectByRecycle[typeName]; ok {
		return point
	}
	return 0
}

func (srv RecycleService) GetType(typeName string) entity.PointTransactionType {
	for name, typeN := range pointRecycleByRules {
		if typeName == name || strings.Contains(name, typeName) {
			return typeN
		}
	}
	return ""
}

func (srv RecycleService) getPoint(name string, number int) (point, co2 int64) {
	if name == "" || number == 0 {
		return point, co2
	}
	//小分类 point / 每台
	if pointByOne, ok := pointRecycleByQuantity[name]; ok {
		point = pointByOne * int64(number)
	}
	//大分类 按公斤算
	if name == "衣帽鞋包" {
		return srv.checkWeightByClothing(int64(number))
	}
	if name == "书籍课本" {
		return srv.checkWeightByBook(int64(number))
	}
	return point, co2
}

func (srv RecycleService) checkWeightByClothing(weight int64) (int64, int64) {
	switch {
	case weight >= 5 && weight <= 10:
		return 208, 91
	case weight >= 11 && weight <= 20:
		return 417, 181
	case weight >= 21 && weight <= 65:
		return 834, 362
	case weight > 65:
		return 2709, 1178
	default:
		return 0, 0
	}
}

func (srv RecycleService) checkWeightByBook(weight int64) (int64, int64) {
	switch {
	case weight >= 10 && weight <= 20:
		return 64, 1
	case weight >= 21 && weight <= 30:
		return 129, 3
	case weight >= 31 && weight <= 50:
		return 193, 4
	case weight > 50:
		return 322, 7
	default:
		return 0, 0
	}
}
