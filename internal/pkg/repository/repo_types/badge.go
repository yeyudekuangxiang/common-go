package repo_types

type FindBadgeBy struct {
	ID      int64
	OrderId string
}

type GetBadgeListBy struct {
	OpenId string
}
