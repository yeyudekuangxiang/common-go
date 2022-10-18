package srv_types

import "github.com/aliyun/alibaba-cloud-sdk-go/services/sts"

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
	OssPath      string
	OperatorId   int64
	OperatorType int8 //1用户 2管理员
	SceneId      int
}

type UploadCallbackParam struct {
	LogId    string
	Filename string
	Size     int64
	MimeType string
	Height   float64
	Width    float64
}

type OssStsInfo struct {
	Credentials sts.Credentials `json:"credentials"`
	Region      string          `json:"region"`
	Bucket      string          `json:"bucket"`
	MimeTypes   []string        `json:"mimeTypes"`
	MaxSize     int64           `json:"maxSize"`
	UploadId    string          `json:"uploadId"`
	Path        string          `json:"path"`
	MaxAge      int             `json:"maxAge"`
}
