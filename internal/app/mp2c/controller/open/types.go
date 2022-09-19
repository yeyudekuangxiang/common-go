package open

type platform struct {
	PlatformKey string `json:"platformKey" form:"platformKey"`
	MemberId    int64  `json:"memberId,omitempty" form:"memberId"`
	Method      string `json:"method" form:"method"`
	Mobile      string `json:"mobile" form:"mobile"`
	Sign        string `json:"sign" form:"sign"`
}
