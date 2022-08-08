package entity

import "time"

type UserFriend struct {
	ID        int
	Uid       int64     `json:"uid"`
	FUid      int64     `json:"fUid"`
	Type      int64     `json:"type"`
	Source    int64     `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
