package entity

type BdSceneUser struct {
	ID          int    `json:"ID"`
	Ch          string `json:"ch"`
	SceneUserId int64  `json:"sceneUserId,omitempty"` //外站用户id
	Phone       string `json:"phone,omitempty"`       //外站用户手机
	OpenId      string `json:"openId,omitempty"`      //本站用户openId
}

func (BdSceneUser) TableName() string {
	return "bd_scene_user"
}
