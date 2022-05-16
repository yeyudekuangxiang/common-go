package entity

import "mio/internal/pkg/model"

type Event struct {
	ID                    int64
	EventCateGoryId       string
	EventId               string
	Title                 string
	Subtitle              string
	Active                bool
	CoverImageUrl         string
	StartTime             model.Time
	EndTime               model.Time
	ProductItemId         string
	ParticipationCount    int
	ParticipationTitle    string
	ParticipationSubtitle string
}
