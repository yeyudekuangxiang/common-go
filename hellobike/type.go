package hellobike

type SendCouponParam struct {
	AppId        string `json:"app_id"`
	Action       string `json:"action"`
	UtcTimestamp string `json:"utc_timestamp"`
	Sign         string `json:"sign"`
	Version      string `json:"version"`
	BizContent   struct {
		TransactionId string `json:"transactionId"`
		MobilePhone   string `json:"mobilePhone"`
		ActivityId    string `json:"activityId"`
	} `json:"biz_content"`
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
