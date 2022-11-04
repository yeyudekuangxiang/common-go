package api_types

type GetGetBannerListForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=home event topic user" alias:"banner场景"`
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
