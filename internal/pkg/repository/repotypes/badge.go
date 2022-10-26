package repotypes

type FindBadgeBy struct {
	ID      int64
	OrderId string
	OpenId  string
}

type GetBadgeListBy struct {
	OpenId string
}
