package entity

type BdScene struct {
	ID         int
	Ch         string
	Key        string
	PointLimit int
	Override   int
	WhiteIp    string
}

func (BdScene) TableName() string {
	return "bd_scene"
}
