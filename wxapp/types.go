package wxapp

type Response struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type QRCodeResponse struct {
	Response
	ContentType string `json:"contentType"`
	Buffer      []byte `json:"buffer"`
}

type internalAppId struct {
	AppId string `json:"appid"`
}
type UserRiskRankParam struct {
	internalAppId
	OpenId       string `json:"openid"`
	Scene        int64  `json:"scene"`
	MobileNo     string `json:"mobile_no"`
	ClientIp     string `json:"client_ip"`
	EmailAddress string `json:"email_address"`
	ExtendedInfo string `json:"extended_info"`
	IsTest       bool   `json:"is_test"`
}
type UserRiskRankResponse struct {
	Response
	UnoinId  int `json:"unoin_id"`
	RiskRank int `json:"risk_rank"`
}
type EnvVersion string

const (
	EnvVersionRelease EnvVersion = "release"
	EnvVersionTrial   EnvVersion = "trial"
	EnvVersionDevelop EnvVersion = "develop"
)

type URLSchemeRequest struct {
	// 跳转到的目标小程序信息。
	SchemedInfo *SchemedInfo `json:"jump_wxa,omitempty"`
	// 成的scheme码类型，到期失效：true，永久有效：false。
	IsExpire bool `json:"is_expire,omitempty"`
	// 到期失效的scheme码的失效时间，为Unix时间戳。生成的到期失效scheme码在该时间前有效。最长有效期为1年。生成到期失效的scheme时必填。
	ExpireTime int64 `json:"expire_time,omitempty"`
}

type SchemedInfo struct {
	// 通过scheme码进入的小程序页面路径，必须是已经发布的小程序存在的页面，不可携带query。path为空时会跳转小程序主页。
	Path string `json:"path"`
	// 通过scheme码进入小程序时的query，最大128个字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~
	Query      string     `json:"query"`
	EnvVersion EnvVersion `json:"env_version"`
}
