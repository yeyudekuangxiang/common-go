package api_types

import "mio/internal/pkg/model/entity"

type DuiBaActivityVO struct {
	ID          int                         `json:"id"`
	Name        string                      `json:"name"`
	Cid         int64                       `json:"cid"`
	Type        entity.DuiBaActivityType    `json:"type"`
	IsShare     entity.DuiBaActivityIsShare `json:"isShare"`
	IsPhone     entity.DuiBaActivityIsPhone `json:"isPhone"`
	ActivityUrl string                      `json:"activityUrl"`
	ActivityId  string                      `json:"activityId"`
	CreatedAt   string                      `json:"createdAt"`
	UpdatedAt   string                      `json:"updatedAt"`
	RiskLimit   int                         `json:"riskLimit"` //用户风险等级限制
}

type DuiBaActivityShowVO struct {
	ID            int    `json:"id"`
	NoLoginH5Link string `json:"noLoginH5Link"`
	StaticH5Link  string `json:"staticH5Link"`
	InsideLink    string `json:"insideLink"`
	EwmLink       string `json:"ewmLink"`
}
