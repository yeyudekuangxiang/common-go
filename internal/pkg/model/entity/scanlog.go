package entity

import (
	"mio/internal/pkg/model"
	"time"
)

type ScanLog struct {
	ID         int64
	ImageUrl   string
	Hash       string //sha256
	Count      int
	ScanResult model.ArrayString
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
