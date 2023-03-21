package usermsg

import "encoding/json"

type IRecycleInfo interface {
	IUHDIOUHOIQWHIOEIOWEIOWEASOKA()
	JSON() ([]byte, error)
}

type BindMobile struct {
	UserId    int64  `json:"userId"`
	OpenId    string `json:"openId"`
	ChannelId int64  `json:"channelId"`
}

func (d BindMobile) IUHDIOUHOIQWHIOEIOWEIOWEASOKA() {

}
func (d BindMobile) JSON() ([]byte, error) {
	return json.Marshal(d)
}

type Interaction struct {
	Tp         string `json:"tp"`
	Data       string `json:"data"`
	Ip         string `json:"ip"`
	Result     string `json:"result"`
	ResultCode string `json:"resultCode"`
	UserId     int64  `json:"userId"`
}

func (d Interaction) IUHDIOUHOIQWHIOEIOWEIOWEASOKA() {

}
func (d Interaction) JSON() ([]byte, error) {
	return json.Marshal(d)
}
