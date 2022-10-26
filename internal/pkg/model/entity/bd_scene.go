package entity

type BdScene struct {
	ID            int    `json:"ID"`
	Ch            string `json:"ch"`
	Key           string `json:"key"`
	PointLimit    int    `json:"pointLimit"`
	Override      int    `json:"override"`
	WhiteIp       string `json:"whiteIp"`
	Secret        string `json:"secret"`
	AppId         string `json:"appId"`
	Domain        string `json:"domain"`
	Secret2       string `json:"secret2"`
	AppId2        string `json:"appId2"`
	Domain2       string `json:"domain2"`
	PrePointLimit int    `json:"prePointLimit"`
}

func (BdScene) TableName() string {
	return "bd_scene"
}
