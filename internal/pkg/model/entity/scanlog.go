package entity

import "time"

type ScanLog struct {
	ID        int64
	ImageUrl  string
	Hash      string //sha256
	Count     int
	CreatedAt time.Time
	UpdatedAt time.Time
}
