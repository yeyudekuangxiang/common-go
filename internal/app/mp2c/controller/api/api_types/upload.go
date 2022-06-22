package api_types

type GetUploadTokenInfoForm struct {
	Scene string `json:"scene" form:"scene" binding:"required" alias:"上传场景"`
}
type OssUploadCallbackForm struct {
	Filename string  `json:"filename" form:"filename" binding:"required" alias:"文件名称"`
	Size     int64   `json:"size" form:"size" binding:"required" alias:"文件大小"`
	MimeType string  `json:"mimeType" form:"mimeType" binding:"required" alias:"文件类型"`
	Height   float64 `json:"height" form:"height" binding:"required" alias:"文件高度"`
	Width    float64 `json:"width" form:"width" binding:"required" alias:"文件宽度"`
}
