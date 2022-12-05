package usermsg

type BindMobile struct {
	UserId    int64  `json:"userId"`
	OpenId    string `json:"openId"`
	ChannelId int64  `json:"channelId"`
}
