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

type BaseResponse struct {
	ErrorCode string `json:"error_code"`
	Success   bool   `json:"success"`
	Data      string `json:"data"`
}

func (c *BaseResponse) IsSuccess() bool {
	return c.ErrorCode == "10000"
}

func (c *BaseResponse) GetErrMsg() string {
	errMsg, ok := ErrorCodeList[c.ErrorCode]
	if ok {
		return errMsg
	}
	return c.ErrorCode
}

var ErrorCodeList = map[string]string{
	"10101":  "系统错误",
	"10104":  "参数异常",
	"10114":  "用户异常",
	"10302":  "appId无效",
	"10500":  "sign校验失败",
	"10501":  "权益校验失败",
	"200000": "非法活动id",
	"200001": "非法手机号",
	"200002": "当前活动已结束",
	"200003": "没找到对应的活动",
	"200004": "活动库存不足",
	"200005": "该用户不存在",
	"200006": "卡发放失败，请重试",
	"200007": "重复参加活动",
}

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
