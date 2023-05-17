package hellobike

type SendCouponParam struct {
	ActivityId    string `json:"activityId"`
	MobilePhone   string `json:"mobilePhone"`
	TransactionId string `json:"transactionId"`
}

type SendHelloBikeCouponParam struct {
	AppId        string `json:"app_id"`
	Action       string `json:"action"`
	UtcTimestamp string `json:"utc_timestamp"`
	Sign         string `json:"sign"`
	Version      string `json:"version"`
	BizContent   struct {
		ActivityId    string `json:"activityId"`
		MobilePhone   string `json:"mobilePhone"`
		TransactionId string `json:"transactionId"`
	} `json:"biz_content"`
}

type ResultCode string

const (
	ResultCodeSuccess ResultCode = "0000"
)

type BaseResponse struct {
	ErrorCode string `json:"error_code"`
	Success   bool   `json:"success"`
	Data      string `json:"data"`
}

/*{"orderNo":"TP20230317180734200102903590028"}*/
/**
{
    "code": 302,
    "msg": "appId无效"
}
*/

type RefundCardResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type RefundCardParam struct {
	ActivityId    string `json:"activityId"`
	OrderNo       string `json:"orderNo"`
	MobilePhone   string `json:"mobilePhone"`
	TransactionId string `json:"transactionId"`
}
type RefundHelloBikeCardParam struct {
	AppId      string `json:"appId"`
	Action     string `json:"action"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
	Version    string `json:"version"`
	BizContent struct {
		ActivityId    string `json:"activityId"`
		OrderNo       string `json:"orderNo"`
		MobilePhone   string `json:"mobilePhone"`
		TransactionId string `json:"transactionId"`
	} `json:"bizContent"`
}
