package open

type platform struct {
	PlatformKey string `json:"platformKey" form:"platformKey"`
	MemberId    int64  `json:"memberId" form:"memberId"`
}
