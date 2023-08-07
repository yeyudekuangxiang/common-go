package point

import (
	"fmt"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
)

//检查图片
func (c *DefaultClientHandle) scanImage(imgUrl string) ([]string, error) {
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
func (c *DefaultClientHandle) validateRule(content []string, rules CollectRules) ([]string, error) {
	ruleArray := util.Intersect(content, rules[c.clientHandle.Type])
	if len(ruleArray) == 0 {
		return nil, errno.ErrCommon.WithMessage("不是有效的图片")
	}
	return content, nil
}

//图片识别 根据rule匹配关键数据
func (c *DefaultClientHandle) identifyImg(identify []string) (map[string]string, error) {
	rule := identifyChRules[c.clientHandle.Type] //汉字
	enRules := identifyEnRules[c.clientHandle.Type]
	m, valid := util.IntersectContains(identify, rule)
	if !valid {
		return nil, errno.ErrCommon.WithMessage("无效图片")
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
func (c *DefaultClientHandle) getText() string {
	return service.PointTransactionTypeInfo{Type: c.clientHandle.Type}.Type.RealText()
}

//平台渠道
func (c *DefaultClientHandle) getRealText() string {
	return service.PointTransactionTypeInfo{Type: c.clientHandle.Type}.Type.Text()
}

//写日志
func (c *DefaultClientHandle) writeMessage(code int, message string) {
	app.Logger.Info(fmt.Sprintf("%d-%s", code, message))
}
