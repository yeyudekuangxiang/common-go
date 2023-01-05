package kumiaoCommunity

import "mio/internal/pkg/model/entity"

type GetActivitiesTagPageListParams struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Offset      int                `json:"offset"`
	Limit       int                `json:"limit"` //limit为0时不限制数量
	OrderBy     entity.OrderByList `json:"orderBy"`
}

type GetActivitiesTagListParams struct {
}
