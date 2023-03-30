package api_types

type GetGetBannerListForm struct {
	Scene   string `json:"scene" form:"scene" binding:"oneof=home event topic user welfare" alias:"banner场景"`
	Display string `json:"display" form:"display" binding:""`
}

type BannerVO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
	Type     string `json:"type"`
	Redirect string `json:"redirect"`
	AppId    string `json:"appId"`
	Ext      string `json:"ext"`
}
