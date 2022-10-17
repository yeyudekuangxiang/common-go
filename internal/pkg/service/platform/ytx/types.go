package ytx

type SynchroRequest struct {
	OpenId         string `json:"openId"`              //亿通行openId
	RegDate        string `json:"regDate"`             //注册时间，格式yyyyMMddHHmmss
	PlatformUserId string `json:"platformUserId"`      //绿喵用户ID
	Ts             int64  `json:"ts"`                  //时间戳，毫秒
	Signature      string `json:"signature,omitempty"` //签名，计算获得
}

type synchroResponse struct {
	ResCode    string                 `json:"resCode"`    //返回码
	ResMessage string                 `json:"resMessage"` //返回描述
	ResData    map[string]interface{} `json:"resData"`    //object
}
