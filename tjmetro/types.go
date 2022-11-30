package tjmetro

type TicketAllotParam struct {
	//发放id 由天津地铁提供
	AllotId string `json:"allotId"`
	//用户id 用户id和手机号二选一
	EtUserId string `json:"etUserId"`
	//用户手机号
	EtUserPhone string `json:"etUserPhone"`
	//发放数量
	AllotNum int `json:"allotNum"`
}

type ResultCode string

const (
	ResultCodeSuccess ResultCode = "0000"
)

type BaseResponse struct {
	//0000表示成功 其他表示失败
	ResultCode ResultCode
	ResultDesc string
	ResultData interface{}
}
