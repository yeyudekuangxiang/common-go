package repotypes

type GetUserFriendListBy struct {
	Uid      int64
	FUid     int64
	Type     int
	UserIds  []int64
	FUserIds []int64
}
