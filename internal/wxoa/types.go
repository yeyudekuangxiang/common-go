package wxoa

type Error struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
type AccessTokenResponse struct {
	Error
	AccessToken
}
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type TicketResponse struct {
	Error
	Ticket
}
type Ticket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}
