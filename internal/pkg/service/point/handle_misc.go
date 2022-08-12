package point

import (
	"errors"
	"fmt"
	"math"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
)

//检查图片
func (c *defaultClientHandle) scanImage(imgUrl string) ([]string, error) {
	results, err := service.DefaultOCRService.Scan(imgUrl)
	if err != nil {
		return nil, err
	}
	ruleArray, err := c.validateRule(results, rules)
	if err != nil {
		return nil, err
	}
	return ruleArray, nil
}

//规则校验
func (c *defaultClientHandle) validateRule(content []string, rules CollectRules) ([]string, error) {
	ruleArray := util.Intersect(content, rules[c.clientHandle.Type])
	if len(ruleArray) == 0 {
		return nil, errors.New("不是有效的图片")
	}
	return ruleArray, nil
}

//图片识别
func (c *defaultClientHandle) identifyImg(identify []string) {
	for i, str := range identify {
		fmt.Printf("%d-%s", i, str)
	}
	return
}

//积分类型
func (c *defaultClientHandle) getText() string {
	text, ok := commandText[c.clientHandle.Type]
	if !ok {
		return "未知积分"
	}
	return text
}

//平台渠道
func (c *defaultClientHandle) getRealText() string {
	text, ok := commandRealText[c.clientHandle.Type]
	if !ok {
		return "未知平台"
	}
	return text
}

//诸葛埋点
func (c *defaultClientHandle) trackPoint(err error) {
	var isFail bool
	if err != nil {
		isFail = true
	}
	c.plugin.tracking.TrackPoints(srv_types.TrackPoints{
		OpenId:      c.clientHandle.OpenId,
		PointType:   c.getRealText(),
		ChangeType:  c.additional.changeType,
		Value:       uint(math.Abs(float64(c.clientHandle.point))),
		IsFail:      isFail,
		FailMessage: err.Error(),
	})
}

//写日志
func (c *defaultClientHandle) writeMessage(code int, message string) {
	app.Logger.Info(fmt.Sprintf("%d-%s", code, message))
}
