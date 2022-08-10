package point

import (
	"errors"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
)

func (c *clientHandle) scanImage(imgUrl string) ([]string, error) {
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

func (c *clientHandle) validateRule(content []string, rules CollectRules) ([]string, error) {
	ruleArray := util.Intersect(content, rules[c.Type])
	if len(ruleArray) == 0 {
		return nil, errors.New("未匹配到对应规则")
	}
	return ruleArray, nil
}
