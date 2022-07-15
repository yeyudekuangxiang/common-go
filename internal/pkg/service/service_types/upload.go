package service_types

type UploadTokenInfo struct {
	OssPolicyToken OssPolicyToken `json:"ossPolicyToken"`
	MimeTypes      []string       `json:"mimeTypes"`
	MaxSize        int64          `json:"maxSize"`
	UploadId       string         `json:"uploadId"`
	Domain         string         `json:"domain"`
	MaxAge         int            `json:"maxAge"`
}
type FindSceneParam struct {
	Scene string
}

type CreateUploadLogParam struct {
	OssPath string
	UserId  int64
	SceneId int
}

type UploadCallbackParam struct {
	LogId    string
	Filename string
	Size     int64
	MimeType string
	Height   float64
	Width    float64
}
