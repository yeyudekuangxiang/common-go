package entity

type BdScene struct {
	ID         int
	Ch         string
	Key        string
	PointLimit int
	Override   int
	WhiteIp    string
	AppId      string
	Domain     string
}

func (BdScene) TableName() string {
	return "bd_scene"
}
