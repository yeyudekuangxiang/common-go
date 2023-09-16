package carbonmsg

type CarbonChangeSuccess struct {
	Openid        string  `json:"openid"`
	UserId        int64   `json:"userId"`
	TransactionId string  `json:"transactionId"`
	Type          string  `json:"type"`
	City          string  `json:"city"`
	Value         float64 `json:"value"`
	Info          string  `json:"info"`
}
