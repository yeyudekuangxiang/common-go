package badge

type UpdateBadgeImageForm struct {
	UploadCode string `json:"uploadCode" form:"uploadCode" binding:"required" alias:"上传凭证"`
	ImageUrl   string `json:"imageUrl" form:"imageUrl" binding:"required" alias:"证书图片"`
}
type UpdateBadgeIsNewForm struct {
	BadgeId int64 `json:"badgeId" form:"badgeId" binding:"required" alias:"证书编号"`
}
