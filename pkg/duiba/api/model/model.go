package model

type Base struct {
	AppKey    string `json:"appKey" form:"appKey" binding:"required" alias:"appKey"`          //接口appKey，应用的唯一标识
	Timestamp IntStr `json:"timestamp" form:"timestamp" binding:"required" alias:"timestamp"` //1970-01-01开始的时间戳，毫秒为单位。
	Uid       string `json:"uid" form:"uid" binding:"required" alias:"uid"`                   //用户唯一性标识
	Sign      string `json:"sign" form:"sign"  alias:"sign"`
}

func (b Base) ToMap() map[string]string {
	return map[string]string{
		"appKey":    b.AppKey,
		"timestamp": string(b.Timestamp),
		"uid":       b.Uid,
		"sign":      b.Sign,
	}
}
