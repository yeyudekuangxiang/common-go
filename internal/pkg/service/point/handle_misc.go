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
	ocrSrv := service.DefaultOCRService()
	imageHash, err := ocrSrv.CheckImageScanCount(imgUrl, 1)
	if err != nil {
		return nil, err
	}

	results, err := ocrSrv.ScanWithHash(imgUrl, imageHash)
	if err != nil {
		return nil, err
	}
	content, err := c.validateRule(results, rules)
	if err != nil {
		return nil, err
	}
	return content, nil
}

//规则校验
func (c *defaultClientHandle) validateRule(content []string, rules CollectRules) ([]string, error) {
	ruleArray := util.Intersect(content, rules[c.clientHandle.Type])
	if len(ruleArray) == 0 {
		return nil, errors.New("不是有效的图片")
	}
	return content, nil
}

//图片识别 根据rule匹配关键数据
func (c *defaultClientHandle) identifyImg(identify []string) (map[string]string, error) {
	rule := identifyChRules[c.clientHandle.Type] //汉字
	enRules := identifyEnRules[c.clientHandle.Type]
	m, valid := util.IntersectContains(identify, rule)
	if !valid {
		return nil, errors.New("无效图片")
	}
	enM := map[string]string{}
	for k, v := range m {
		if enRule, ok := enRules[k]; ok {
			enM[enRule] = v
		}
	}
	c.clientHandle.identifyImg = enM
	return enM, nil
}

//积分类型
func (c *defaultClientHandle) getText() string {
	return service.PointTransactionTypeInfo{Type: c.clientHandle.Type}.Type.RealText()
}

//平台渠道
func (c *defaultClientHandle) getRealText() string {
	return service.PointTransactionTypeInfo{Type: c.clientHandle.Type}.Type.Text()
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
