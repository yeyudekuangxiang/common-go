package admtypes

type GetUploadTokenInfoForm struct {
	Scene string `json:"scene" form:"scene" binding:"required" alias:"上传场景"`
}
