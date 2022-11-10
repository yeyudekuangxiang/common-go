package api_types

import "mio/internal/pkg/model/entity"

type IndexIconVO struct {
	ID        int64                  `json:"id"`
	Title     string                 `json:"title"`
	Type      string                 `json:"type"`
	RowNum    string                 `json:"rowNum"`
	Sort      int8                   `json:"sort"`
	Status    entity.IndexIconStatus `json:"status"`
	IsOpen    int8                   `json:"isOpen"`
	Pic       string                 `json:"pic"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
}
