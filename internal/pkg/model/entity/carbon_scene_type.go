package entity

import "time"

type CarbonSceneType struct {
	ID        int64               `json:"id"`
	Type      CarbonSceneTypeType `json:"type"`
	Title     string              `json:"title"`
	Icon      string              `json:"icon"`
	Desc      string              `json:"desc"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
}

type CarbonSceneTypeType string
