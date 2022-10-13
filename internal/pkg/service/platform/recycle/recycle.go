package recycle

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
	"mio/pkg/errno"
	"sort"
	"strconv"
	"time"
)

func NewRecycleService(ctx *context.MioContext) *RecycleService {
	return &RecycleService{
		ctx: ctx,
	}
}

type RecycleService struct {
	ctx *context.MioContext
}

//每个大分类对应的每月获取积分上限
var recycleMonthPoint = map[entity.PointTransactionType]int64{
	entity.POINT_RECYCLING_CLOTHING:     2709, // 旧衣回收
	entity.POINT_RECYCLING_BOOK:         322,  // 书籍课本
	entity.POINT_RECYCLING_DIGITAL:      1911, // 数码产品
	entity.POINT_RECYCLING_APPLIANCE:    1840, // 家电回收
	entity.POINT_FMY_RECYCLING_CLOTHING: 2709,
}

//获取分类
var pointRecycleByRules = map[string]entity.PointTransactionType{
	"衣帽鞋包": entity.POINT_RECYCLING_CLOTHING,  // 人/月
	"书籍课本": entity.POINT_RECYCLING_BOOK,      // 人/月
	"数码产品": entity.POINT_RECYCLING_DIGITAL,   // 人/月
	"家电回收": entity.POINT_RECYCLING_APPLIANCE, // 人/月
	"旧物回收": entity.POINT_FMY_RECYCLING_CLOTHING,
}

// 回收 台 单位对应积分 比如 电视机 1台 获得 69积分
var recyclePointByNum = map[string]int64{
	"手机":    113,
	"平板电脑":  409,
	"手提电脑":  1031,
	"一体机电脑": 1911,
	"冰箱":    384,
	"洗衣机":   690,
	"空调":    205,
	"电视机":   69,
	"衣帽鞋包":  21, //1000g : 21 积分
	"书籍课本":  6,  //1000g : 6 积分
	"旧物回收":  21,
}

// 回收 台/重量 单位对应减碳量 比如 电视机 1台 获得 15000g 减碳量
var recycleCo2ByNum = map[string]float64{
	"手机":    25000,
	"平板电脑":  89000,
	"手提电脑":  224000,
	"一体机电脑": 415000,
	"冰箱":    83000,
	"洗衣机":   150000,
	"空调":    45000,
	"电视机":   15000,
	"衣帽鞋包":  4500, //1000g : 4500g
	"书籍课本":  1400, //1000g : 1400g
	"旧物回收":  4500,
}

//每个类型对应次数
var recycleLimit = map[string]int{
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
	"旧物回收":  1,
}

// CheckLimit 检查该类型今日获取次数
func (srv RecycleService) CheckLimit(openId, typeName string) error {
	if err := srv.checkLimit(openId, typeName); err != nil {
		return err
	}
	return nil
}

func (srv RecycleService) CheckOrder(openId, OrderNo string) error {
	if err := srv.checkOrder(openId, OrderNo); err != nil {
		return err
	}
	return nil
}

// CheckSign 验证签名
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
		signStr += v + "=" + util.InterfaceToString(params[v]) + ";"
	}
	//验证签名
	signMd5 := encrypt.Md5(secret + signStr)
	if signMd5 != sign {
		return errno.ErrValidation
	}
	return nil
}

func (srv RecycleService) CheckFmySign(params FmySignParams, appId string, secret string) error {
	sign := params.Sign
	params.Sign = ""
	rand1 := string([]rune(sign)[0:4])
	rand2 := string([]rune(sign)[len(sign)-4:])
	marshal, _ := json.Marshal(params)
	fmt.Printf("marshal:%s\n", marshal)
	verifyData := rand1 + appId + string(marshal) + secret + rand2
	md5Str := encrypt.Md5(verifyData)
	signMd5 := rand1 + md5Str[7:28] + rand2
	//验证签名
	if signMd5 != sign {
		return errno.ErrValidation
	}
	return nil
}

// GetMaxPointByMonth 获取该类型本月积分上限
func (srv RecycleService) GetMaxPointByMonth(typeName entity.PointTransactionType) (int64, error) {
	point := srv.getMaxPointByMonth(typeName)
	if point == 0 {
		return 0, errors.New("")
	}
	return point, nil
}

// GetType 获取大类型
func (srv RecycleService) GetType(typeName string) entity.PointTransactionType {
	for name, typeN := range pointRecycleByRules {
		if typeName == name || typeName == string(typeN) {
			return typeN
		}
	}
	return ""
}

func (srv RecycleService) GetText(typeName entity.PointTransactionType) string {
	for name, typeN := range pointRecycleByRules {
		if typeName == typeN {
			return name
		}
	}
	return ""
}

// GetPoint 获取积分 qua: 数量&重量 unit: 公斤，个
func (srv RecycleService) GetPoint(typeName, qua string) (int64, error) {
	num, _ := strconv.ParseFloat(qua, 64)
	point := srv.getPoint(typeName, num)
	if point == 0 {
		return 0, errors.New("未匹配到" + typeName + "积分规则")
	}
	return point, nil
}

func (srv RecycleService) GetCo2(typeName, qua string) (float64, error) {
	num, _ := strconv.ParseFloat(qua, 64)
	co2 := srv.getCo2(typeName, num)
	if co2 == 0 {
		return 0, errors.New("未匹配到" + typeName + "减碳量规则")
	}
	return co2, nil
}

func (srv RecycleService) checkLimit(openId string, typeName string) error {
	num := recycleLimit[typeName]
	timeStr := time.Now().Format("20060102")
	limit, _ := app.Redis.Get(srv.ctx, openId+typeName+timeStr).Int()
	if limit >= num {
		return errors.New(typeName + "回收积分已达到本日次数上限")
	}
	app.Redis.Set(srv.ctx, openId+typeName+timeStr, num, time.Hour*24)
	return nil
}

func (srv RecycleService) getMaxPointByMonth(typeName entity.PointTransactionType) int64 {
	if point, ok := recycleMonthPoint[typeName]; ok {
		return point
	}
	return 0
}

//返回积分
func (srv RecycleService) getPoint(typeName string, number float64) int64 {
	var point int64
	if typeName == "" || number == 0 {
		return 0
	}
	//获取point
	if pointByOne, ok := recyclePointByNum[typeName]; ok {
		point = pointByOne * int64(math.Floor(number))
	}
	return point
}

//返回减碳量 单位 g
func (srv RecycleService) getCo2(typeName string, number float64) float64 {
	var co2 float64
	if typeName == "" || number == 0 {
		return 0
	}
	//获取co2
	if co2ByOne, ok := recycleCo2ByNum[typeName]; ok {
		co2 = co2ByOne * math.Floor(number)
	}
	return co2
}

func (srv RecycleService) checkOrder(openId, orderNo string) error {
	pointTransactionRepository := repository.NewPointTransactionRepository(srv.ctx)
	one := pointTransactionRepository.FindBy(repository.FindPointTransactionBy{
		OpenId: openId,
		Note:   orderNo,
	})
	if one.ID != 0 {
		return errors.New("重复订单")
	}
	return nil
}

func (srv RecycleService) getText(typeText entity.PointTransactionType) string {
	switch typeText {
	case entity.POINT_RECYCLING_CLOTHING:
		return "衣物鞋帽"
	case entity.POINT_RECYCLING_DIGITAL:
		return "数码产品"
	case entity.POINT_RECYCLING_BOOK:
		return "书籍课本"
	case entity.POINT_RECYCLING_APPLIANCE:
		return "家电回收"
	case entity.POINT_FMY_RECYCLING_CLOTHING:
		return "旧物回收"
	default:
		return "未知类型"
	}
}

func (srv RecycleService) getPointType(typeText string) entity.PointTransactionType {
	switch typeText {
	case "衣物鞋帽":
		return entity.POINT_RECYCLING_CLOTHING
	case "数码产品":
		return entity.POINT_RECYCLING_DIGITAL
	case "书籍课本":
		return entity.POINT_RECYCLING_BOOK
	case "家电回收":
		return entity.POINT_RECYCLING_APPLIANCE
	case "旧物回收":
		return entity.POINT_FMY_RECYCLING_CLOTHING
	default:
		return entity.PointTransactionType("未知类型")
	}
}
