package platform

//{"errCode":"0001","message":"今日已打卡，无需重复打卡"}
type zyhCommonResponse struct {
	Data    string `json:"data"`
	ErrCode string `json:"errCode"`
	Message string `json:"message"`
}
