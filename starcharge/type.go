package starcharge

type SendStarChargeParam struct {
	Data     string `json:"data"`
	QueryUrl string `json:"queryUrl"`
}
type starChargeResponse struct {
	Ret  int    `json:"Ret"`
	Msg  string `json:"Msg"`
	Data string `json:"Data"`
	Sig  string `json:"Sig"`
}
