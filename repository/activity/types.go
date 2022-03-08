package activity

type FindRecordBy struct {
	UserId int64
}
type GetRecordListBy struct {
	ApplyStatus int8 `binding:"oneof=0 1 2 3 4" alias:"申请状态"` //0 全部 1参与中 2申请中 3申请成功 4申请失败
	UserIds     []int64
	Offset      int
	Limit       int
}
